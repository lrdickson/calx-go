package controller

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func NameValid(input string) error {
	// Check for valid characters
	letters := `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	letters += `abcdefghijklmnopqrstuvwxyz`
	validCharacters := letters
	validCharacters += `0123456789`
	validCharacters += `_`
	for index, character := range input {
		characterString := string(character)
		if index == 1 && !strings.Contains(letters, characterString) {
			log.Println("Invalid variable name")
			return errors.New(`"` + characterString + "\" is not a valid 1st character")
		}
		if !strings.Contains(validCharacters, characterString) {
			log.Println("Invalid variable name")
			return errors.New(`"` + characterString + "\" is not a valid character")
		}
	}
	return nil
}

type ObjectEngine interface {
	Close() error
	Kernel() string
}

type Consumer interface {
	ObjectEngine
	Consume(any) error
}

type BaseObjectEngine struct {
}

func (b *BaseObjectEngine) Close() error { return nil }

func (b *BaseObjectEngine) Kernel() string { return "" }

type Producer interface {
	ObjectEngine
	SetOnOutputChanged(func())
	Output() (any, error)
}

type Object struct {
	controller   *Controller
	dependencies []*Object
	dependents   []*Object
	engine       *ObjectEngine
	name         string
}

func (o *Object) Name() string {
	return o.name
}

func (o *Object) SetName(name string) error {
	// Make sure the name is unique
	if _, exists := o.controller.objectNames[name]; exists {
		return fmt.Errorf("The name %s is taken", name)
	}

	// Make sure the name is valid
	if err := NameValid(name); err != nil {
		return err
	}

	// Update the name
	o.controller.objectNames[name] = o
	delete(o.controller.objectNames, o.name)
	o.name = name
	o.controller.EventTriggered(RenameObjectEvent, o)

	// Success
	return nil
}

func (o *Object) Engine() *ObjectEngine {
	return o.engine
}
