package runner

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

type workerStatus int

const (
	ok workerStatus = iota
	failed
)

type worker struct {
	quit   chan int
	in     chan []any
	out    chan any
	status chan workerStatus
	name   chan string
	active bool
}

func (w *worker) stop() {
	if w.active {
		w.quit <- 0
		w.active = false
	}
}

type Kernel struct {
	workers map[string]*worker
}

func NewKernel() Kernel {
	k := Kernel{workers: make(map[string]*worker)}
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
	for name, code := range workerFormulas {
		// Stop the worker if it already existed
		if oldWorker, exists := k.workers[name]; exists {
			oldWorker.stop()
		}

		// Create a new worker
		quit := make(chan int)
		in := make(chan []any)
		out := make(chan any)
		status := make(chan workerStatus)
		nameChannel := make(chan string)
		newWorker := worker{
			quit:   quit,
			in:     in,
			out:    out,
			status: status,
			name:   nameChannel,
			active: true,
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
				status <- failed
				return
			}

			// Build the function code
			functionCode := "package run\nfunc Run() any {\n"
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
				status <- failed
				return
			}
			v, err := gointerp.Eval("run.Run")
			if err != nil {
				log.Println("Failed to get", name, "function:", err)
				status <- failed
				return
			}
			function := v.Interface().(func() any)

			for {
				select {
				case <-quit:
					fmt.Println("Quiting:", name)
					status <- ok
					return
				case <-in:
					// Get the function output
					out <- function()
				case name = <-nameChannel:
					continue
				}
			}
		}(name, code)
	}

	// Get the output
	outputData := make(map[string]string)
	for _, activeWorker := range k.workers {
		activeWorker.in <- make([]any, 0)
	}
	for name, activeWorker := range k.workers {
		select {
		case <-activeWorker.status:
			continue
		case output := <-activeWorker.out:
			fmt.Println("Sending quit signal to:", name)
			activeWorker.stop()
			fmt.Println("quit signal sent to:", name)
			switch output.(type) {
			case bool:
				outputData[name] = strconv.FormatBool(output.(bool))
			case int:
				outputData[name] = strconv.Itoa(output.(int))
			case uint:
				outputData[name] = strconv.FormatUint(uint64(output.(uint)), 10)
			case float32:
				outputData[name] = strconv.FormatFloat(float64(output.(float32)), 'f', -1, 32)
			case float64:
				outputData[name] = strconv.FormatFloat(output.(float64), 'f', -1, 64)
			case string:
				outputData[name] = output.(string)
			default:
				outputReflect := reflect.ValueOf(output)
				outputData[name] = fmt.Sprint(outputReflect.Kind())
			}
		}
	}
	return outputData
}
