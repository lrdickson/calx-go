package viewmodels

import (
	"reflect"

	"fyne.io/fyne/v2/data/binding"
)

type MainViewModel struct {
	// Public
	VariableList binding.StringList
	ValueDisplay binding.String
	EditorCode   binding.String

	//Private
	variableCode  map[string]string
	variableValue map[string]reflect.Value
}

func NewMainViewModel() MainViewModel {
	return MainViewModel{
		binding.NewStringList(),
		binding.NewString(),
		binding.NewString(),
		make(map[string]string),
		make(map[string]reflect.Value)}
}
