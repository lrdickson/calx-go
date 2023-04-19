package gokernel

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"reflect"
	"sync"

	"github.com/lrdickson/calx/internal/controller"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// TODO: Move kernel worker logic into the formula
type FormulaEngine struct {
	controller.BaseObjectEngine
	code            string
	output          any
	err             error
	onOutputChanged func()
	quit            chan int
	runSignal       chan int
	wait            sync.WaitGroup
}

func NewFormula(c *controller.Controller) *controller.Object {
	// Make the formula object
	formulaEngine := FormulaEngine{}
	var objectEngine controller.ObjectEngine = &formulaEngine
	formula := c.NewObject(c.UniqueName(), &objectEngine)
	log.Println("Creating new formula:", formula.Name())
	quit := make(chan int)
	run := make(chan int)
	formulaEngine.quit = quit
	formulaEngine.runSignal = run

	// Start the new worker
	go func() {
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

		for {
			log.Println(formula.Name(), "ready to receive commands")
			select {
			case <-quit:
				log.Println("Quiting:", formula.Name())
				return
			case <-run:
				// Build the function code
				functionCode := "package run\n"
				functionCode += `import . "math"` + "\n"
				functionCode += "func Run(params []any) any {\n"

				// Get the function parameters
				params := make([]any, 0, len(formula.Dependencies()))
				success := true
				for index, dependency := range formula.Dependencies() {
					// Get the result
					dependency.Wait()
					output, err := dependency.Output()
					if err != nil {
						formulaEngine.err = fmt.Errorf(
							"Failed to run due to input error from: %s", dependency.Name())
						success = false
					}
					params = append(params, output)

					// Determine the type
					if output == nil {
						functionCode += fmt.Sprintf("var %s any = nil", dependency.Name())
					} else {
						paramType := reflect.TypeOf(output)
						functionCode += fmt.Sprintf(
							"%s.(%s) := %d", dependency.Name(), paramType.String(), index)
					}
					functionCode += "\n"
				}
				if !success {
					break
				}

				// Add in the function code
				functionCode += formulaEngine.Code()
				functionCode += "}"

				// Create the function
				log.Println("Function code:\n", functionCode)
				_, err := gointerp.Eval(functionCode)
				if err != nil {
					// TODO: Display this error to the user
					formulaEngine.err = fmt.Errorf(
						"Failed to evaluate %s code: %s", formula.Name(), err.Error())
					log.Println(formulaEngine.err)
					return
				}
				v, err := gointerp.Eval("run.Run")
				if err != nil {
					log.Println("Failed to get", newWorker.name, "function:", err)
					k.status <- workerStatus{newWorker.name, failed}
					return
				}
				function := v.Interface().(func([]any) any)

				// Get the function output
				log.Println(newWorker.name, "running function")
				newWorker.result = func() (result any) {
					defer func() {
						if r := recover(); r != nil {
							// TODO: Display this error to the user
							log.Println("Recoverd from yaegi panic:", r)
							result = r
							return
						}
					}()
					return function(params)
				}()
				log.Println(newWorker.name, "function returned result", newWorker.result)
				newWorker.wait.Done()
				done <- newWorker.name
			}
		}
	}()
	return formula
}

func (f *FormulaEngine) Code() string {
	return f.code
}

func (f *FormulaEngine) SetCode(code string) {
	f.code = code
}

func (f *FormulaEngine) Consume(any) error {
	return nil
}

func (f *FormulaEngine) Output() (any, error) {
	return f.output, nil
}

func (f *FormulaEngine) SetOnOutputChanged(changed func()) {
	f.onOutputChanged = changed
}

func (f *FormulaEngine) Wait() {
	f.wait.Wait()
}
