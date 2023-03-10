package viewmodels

import (
	"log"
	"reflect"

	"github.com/lrdickson/ssgo/internal/runner"
)

type MainViewModel struct {
	// Public
	VariableList []string
	EditorCode   string

	//Private
	variableCode  map[string]string
	variableValue map[string]reflect.Value
}

func NewMainViewModel() MainViewModel {
	return MainViewModel{
		[]string{},
		"",
		make(map[string]string),
		make(map[string]reflect.Value)}
}

func (vm MainViewModel) RunCode() string {
	log.Println("Running:", vm.EditorCode)
	display, err := runner.RunGo(vm.EditorCode)
	if err != nil {
		display = "Err"
	}
	return display
}
