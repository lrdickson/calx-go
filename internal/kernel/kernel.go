package kernel

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const (
	ok int = iota
	failed
)

type workerStatus struct {
	name  string
	value int
}

type worker struct {
	quit     chan int
	in       chan []any
	name     chan string
	active   bool
	response chan any
}

type workerOutput struct {
	name string
	data any
}

type Kernel struct {
	workers map[string]*worker
	status  chan workerStatus
}

func (k *Kernel) stop(name string) {
	if _, exists := k.workers[name]; exists {
		if k.workers[name].active {
			k.workers[name].quit <- 0
			k.workers[name].active = false
			<-k.status
		}
	}
}

func NewKernel() Kernel {
	k := Kernel{workers: make(map[string]*worker)}
	status := make(chan workerStatus)
	k.status = status
	go func() {
		for {
			ws := <-status
			fmt.Println(ws.name, "quit with status", ws.value)
			k.workers[ws.name].active = false
		}
	}()

	return k
}

func (k *Kernel) getWorker(name string) (*worker, bool) {
	if formulaWorker, exists := k.workers[name]; exists {
		return formulaWorker, formulaWorker.active
	}
	return nil, false
}

func (k *Kernel) RenameFormula(oldName, newName string) {
	if formulaWorker, exists := k.workers[oldName]; exists {
		if formulaWorker.active {
			formulaWorker.name <- newName
		}
		k.workers[newName] = formulaWorker
		delete(k.workers, oldName)
	}
}

func (k *Kernel) Update(workerFormulas map[string]string) map[string]string {
	query := make(chan []string)
	out := make(chan workerOutput)
	for name, code := range workerFormulas {
		// Stop the worker if it already existed
		k.stop(name)

		// Create a new worker
		quit := make(chan int)
		in := make(chan []any)
		nameChannel := make(chan string)
		response := make(chan any)
		newWorker := worker{
			quit:     quit,
			in:       in,
			name:     nameChannel,
			active:   true,
			response: response,
		}
		k.workers[name] = &newWorker

		// Start the new worker
		go func(name, code string) {
			// Start the interpreter
			gointerp := interp.New(interp.Options{
				GoPath: build.Default.GOPATH,
				Env:    os.Environ(),
				//Unrestricted: true,
			})
			if err := gointerp.Use(stdlib.Symbols); err != nil {
				log.Fatal("Stdlib load error:", err)
			}
			if err := gointerp.Use(interp.Symbols); err != nil {
				log.Fatal("Interp symbol load error:", err)
			}

			// Add default imports
			_, err := gointerp.Eval(`import "math"`)
			if err != nil {
				log.Println("Failed to import math:", err)
				k.status <- workerStatus{name, failed}
				return
			}

			// Build the function code
			functionCode := "package run\nfunc Run(parent string, query chan []string, response chan any) any {\n"
			functionCode += "get := func(name string) any {\n"
			functionCode += "query <- []string{parent, name}\n"
			functionCode += "return (<- response)\n}\n"
			functionCode += code
			//codeLines := strings.Split(strings.TrimSpace(code), "\n")
			//for lineNumber, line := range codeLines {
			//if lineNumber == len(codeLines)-1 {
			//if !strings.Contains(line, "return") {
			//functionCode += "return "
			//}
			//}
			//functionCode += line + "\n"
			//}
			functionCode += "}"

			// Create the function
			fmt.Println("Function code:\n", functionCode)
			_, err = gointerp.Eval(functionCode)
			if err != nil {
				log.Println("Failed to evaluate", name, "code:", err)
				k.status <- workerStatus{name, failed}
				return
			}
			v, err := gointerp.Eval("run.Run")
			if err != nil {
				log.Println("Failed to get", name, "function:", err)
				k.status <- workerStatus{name, failed}
				return
			}
			function := v.Interface().(func(string, chan []string, chan any) any)

			for {
				select {
				case <-quit:
					fmt.Println("Quiting:", name)
					k.status <- workerStatus{name, ok}
					return
				case <-in:
					// Get the function output
					result := function(name, query, response)
					out <- workerOutput{name, result}
				case name = <-nameChannel:
					continue
				}
			}
		}(name, code)
	}

	// Get the output
	activeWorkerCount := 0
	for _, activeWorker := range k.workers {
		if activeWorker.active {
			activeWorker.in <- make([]any, 0)
			activeWorkerCount++
		}
	}
	//for name, activeWorker := range k.workers {
	outputData := make(map[string]string)
	responseReceived := make(map[string]bool)
	for {
		select {
		case output := <-out:
			name := output.name
			fmt.Println("Sending quit signal to:", name)
			k.stop(name)
			fmt.Println("quit signal sent to:", name)
			switch output.data.(type) {
			case bool:
				outputData[name] = strconv.FormatBool(output.data.(bool))
			case int:
				outputData[name] = strconv.Itoa(output.data.(int))
			case uint:
				outputData[name] = strconv.FormatUint(uint64(output.data.(uint)), 10)
			case float32:
				outputData[name] = strconv.FormatFloat(float64(output.data.(float32)), 'f', -1, 32)
			case float64:
				outputData[name] = strconv.FormatFloat(output.data.(float64), 'f', -1, 64)
			case string:
				outputData[name] = output.data.(string)
			default:
				outputReflect := reflect.ValueOf(output.data)
				outputData[name] = fmt.Sprint(outputReflect.Kind())
			}
			responseReceived[name] = true
		}
		if len(responseReceived) == activeWorkerCount {
			break
		}
	}
	return outputData
}
