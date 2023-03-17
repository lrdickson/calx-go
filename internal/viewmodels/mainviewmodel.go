package viewmodels

import (
	"log"
	"reflect"

	"github.com/lrdickson/ssgo/internal/runner"
)

type Variable struct {
	Code  string
	Value reflect.Value
}

type MainViewModel struct {
	// Public
	EditorCode string
}

func NewMainViewModel() MainViewModel {
	return MainViewModel{""}
}

func (vm MainViewModel) RunCode() string {
	log.Println("Running:", vm.EditorCode)
	display, err := runner.RunGo(vm.EditorCode)
	if err != nil {
		log.Println("Failed to execute code:", err)
		display = "Err"
	}
	return display
}
