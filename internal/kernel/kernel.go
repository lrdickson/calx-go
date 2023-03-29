package kernel

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"

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
	quit   chan int
	in     chan []any
	name   chan string
	active atomic.Bool
	result any
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
		if k.workers[name].active.Load() {
			k.workers[name].quit <- 0
		}
	}
}

func (k *Kernel) getActiveCount() int {
	activeCount := 0
	for name := range k.workers {
		if k.workers[name].active.Load() {
			activeCount++
		}
	}
	return activeCount
}

func NewKernel() *Kernel {
	// Make the kernel
	status := make(chan workerStatus)
	k := &Kernel{
		workers: make(map[string]*worker),
		status:  status,
	}

	// Watch for errors
	go func() {
		for {
			ws := <-status
			fmt.Println(ws.name, "quit with status", ws.value)
			k.workers[ws.name].active.Store(false)
		}
	}()

	return k
}

func (k *Kernel) addWorker(name, code string, done chan string) {
	// Stop the worker if it already existed
	k.stop(name)

	// Create a new worker
	fmt.Println("Creating new worker:", name)
	quit := make(chan int)
	in := make(chan []any)
	nameChannel := make(chan string)
	newWorker := worker{
		quit: quit,
		in:   in,
		name: nameChannel,
	}
	newWorker.active.Store(true)
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

		// Build the function code
		//functionCode := "package run\nfunc Run(params ...any) any {\n"
		functionCode := "package run\n"
		functionCode += `import "math"` + "\n"
		functionCode += "func Run() any {\n"
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
		_, err := gointerp.Eval(functionCode)
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
		//function := v.Interface().(func(...any) any)
		function := v.Interface().(func() any)

		for {
			fmt.Println(name, "ready to receive commands")
			select {
			case <-quit:
				fmt.Println("Quiting:", name)
				k.status <- workerStatus{name, ok}
				return
			//case params := <-in:
			case <-in:
				// Get the function output
				fmt.Println(name, "running function")
				newWorker.result = function()
				fmt.Println(name, "function has run")
				done <- name
			case name = <-nameChannel:
				continue
			}
		}
	}(name, code)
}

func (k *Kernel) getWorker(name string) (*worker, bool) {
	if formulaWorker, exists := k.workers[name]; exists {
		return formulaWorker, formulaWorker.active.Load()
	}
	return nil, false
}

func (k *Kernel) RenameFormula(oldName, newName string) {
	if formulaWorker, exists := k.workers[oldName]; exists {
		if formulaWorker.active.Load() {
			formulaWorker.name <- newName
		}
		k.workers[newName] = formulaWorker
		delete(k.workers, oldName)
	}
}

func (k *Kernel) Update(workerFormulas map[string]string) map[string]string {
	// Make a worker for each formula provided
	done := make(chan string)
	for name, code := range workerFormulas {
		k.addWorker(name, code, done)
	}

	// Get the output
	for _, activeWorker := range k.workers {
		fmt.Println("Starting:", activeWorker.name)
		inputSent := false
		for {
			if activeWorker.active.Load() {
				select {
				case activeWorker.in <- make([]any, 0):
					fmt.Println("Input sent to:", activeWorker.name)
					inputSent = true
				default:
					//fmt.Println(activeWorker.name, "not ready")
					time.Sleep(time.Millisecond)
				}
			} else {
				break
			}

			// escape if input sent
			if inputSent {
				break
			}
		}
	}
	outputData := make(map[string]string)
	responseReceived := make(map[string]bool)
	for {
		select {
		case name := <-done:
			// Get the worker
			fmt.Println("Sending quit signal to:", name)
			activeWorker, workerExists := k.workers[name]
			if !workerExists {
				break
			}

			// Stop the worker
			// TODO: rename workers and leave them running
			k.stop(name)
			fmt.Println("quit signal sent to:", name)

			// Interpet the data
			switch activeWorker.result.(type) {
			case bool:
				outputData[name] = strconv.FormatBool(activeWorker.result.(bool))
			case int:
				outputData[name] = strconv.Itoa(activeWorker.result.(int))
			case uint:
				outputData[name] = strconv.FormatUint(uint64(activeWorker.result.(uint)), 10)
			case float32:
				outputData[name] = strconv.FormatFloat(float64(activeWorker.result.(float32)), 'f', -1, 32)
			case float64:
				outputData[name] = strconv.FormatFloat(activeWorker.result.(float64), 'f', -1, 64)
			case string:
				outputData[name] = activeWorker.result.(string)
			default:
				outputReflect := reflect.ValueOf(activeWorker.result)
				outputData[name] = fmt.Sprint(outputReflect.Kind())
			}
			responseReceived[name] = true
		case <-time.After(time.Millisecond):
			//fmt.Println("timeout")
		}
		if k.getActiveCount() == 0 {
			break
		}
		//fmt.Println("Workers are still active")
	}
	return outputData
}
