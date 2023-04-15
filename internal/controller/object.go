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

type Object interface {
	Close() error
}

type Consumer interface {
	Object
	Consume(any) error
}

type BaseObject struct {
}

func (b *BaseObject) Close() error { return nil }

type Producer interface {
	Object
	SetReady(func())
	Produce() (any, error)
}

type ObjectHolder struct {
	object       *Object
	dependencies []*ObjectHolder
	dependents   []*ObjectHolder
	controller   *Controller
	name         string
}

func (o *ObjectHolder) Name() string {
	return o.name
}

func (o *ObjectHolder) SetName(name string) error {
	// Make sure the name is unique
	if _, exists := o.controller.objectNames[name]; exists {
		return fmt.Errorf("The name %s is taken", name)
	}
	o.controller.objectNames[name] = o
	delete(o.controller.objectNames, o.name)

	// Update the name
	o.name = name
	o.controller.EventTriggered(RenameObjectEvent, o)

	// Success
	return nil
}
