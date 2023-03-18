package runner

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func RunGo(code string) (display string, returnErr error) {
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
		returnErr = fmt.Errorf("Failed to evaluate code: %w", err)
		return
	}

	// Build the function
	codeLines := strings.Split(strings.TrimSpace(code), "\n")
	functionCode := "package run\nfunc Run() any {\n"
	for lineNumber, line := range codeLines {
		if lineNumber == len(codeLines)-1 {
			if !strings.Contains(line, "return") {
				functionCode += "return "
			}
		}
		functionCode += line + "\n"
	}
	functionCode += "}"

	// Run the code
	fmt.Println("Function code:\n", functionCode)
	_, err = gointerp.Eval(functionCode)
	if err != nil {
		returnErr = fmt.Errorf("Failed to evaluate code: %w", err)
		return
	}
	v, err := gointerp.Eval("run.Run")
	if err != nil {
		returnErr = fmt.Errorf("Failed to get function: %w", err)
		return
	}
	function := v.Interface().(func() any)
	output := function()

	// Parse the result
	switch output.(type) {
	case bool:
		display = strconv.FormatBool(output.(bool))
	case int:
		display = strconv.Itoa(output.(int))
	case uint:
		display = strconv.FormatUint(uint64(output.(uint)), 10)
	case float32:
		display = strconv.FormatFloat(float64(output.(float32)), 'f', -1, 32)
	case float64:
		display = strconv.FormatFloat(output.(float64), 'f', -1, 64)
	case string:
		display = output.(string)
	default:
		outputReflect := reflect.ValueOf(output)
		display = fmt.Sprint(outputReflect.Kind())
	}
	return
}
