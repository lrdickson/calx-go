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

	// Run the code
	var v reflect.Value
	for _, line := range strings.Split(code, "\n") {
		var err error
		v, err = gointerp.Eval(line)
		if err != nil {
			returnErr = fmt.Errorf("Failed to evaluate code: %w", err)
			return
		}
	}

	// Parse the result
	switch v.Kind() {
	case reflect.Bool:
		display = strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		display = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		display = strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		display = strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Float64:
		display = strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.String:
		display = v.String()
	default:
		display = "???"
	}
	return
}
