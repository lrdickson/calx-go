package controller

import (
	"log"
	"strconv"

	"github.com/lrdickson/calx/internal/variable"
)

type Event int

const (
	NewVarEvent Event = iota
	RenameVarEvent
	DeleteVarEvent
)

var events []Event = []Event{NewVarEvent, RenameVarEvent, DeleteVarEvent}

// listeners[event][variableName]
type listenerMap map[Event]map[string][]func(string)

type Controller struct {
	variables     map[string]*variable.Variable
	variableCount int
	listeners     listenerMap
}

func NewController() *Controller {
	// Initialize the listeners map
	listeners := make(listenerMap)
	for _, event := range events {
		listeners[event] = make(map[string][]func(string))
		// Add universal listenner
		listeners[event]["*"] = make([]func(string), 0)
	}

	// Create the controller
	return &Controller{
		variables:     make(map[string]*variable.Variable),
		variableCount: 1,
		listeners:     listeners,
	}
}

func (c Controller) IterVariables(iter func(string, *variable.Variable) bool) {
	for key, value := range c.variables {
		cont := iter(key, value)
		if !cont {
			break
		}
	}
}

func (c Controller) Variables(name string) *variable.Variable {
	return c.variables[name]
}

func (c *Controller) uniqueName() string {
	name := ""
	for {
		name = "var" + strconv.Itoa(c.variableCount)
		c.variableCount++
		if _, taken := c.variables[name]; !taken {
			break
		}
	}
	return name
}

func (c *Controller) AddVariable(name string, v *variable.Variable) {
	c.variables[name] = v
	for _, callback := range c.listeners[NewVarEvent]["*"] {
		callback(name)
	}
}

func (c *Controller) AddFormula() {
	var formula variable.Variable = variable.Formula{}
	c.AddVariable(c.uniqueName(), &formula)
}

func (c *Controller) Rename(oldName, newName string) {
	// Check if the oldName exists
	if _, exists := c.variables[oldName]; !exists {
		log.Println("Error: attempt to rename a variable that doesn't exist:", oldName)
		return
	}

	// Update the variable map
	c.variables[newName] = c.variables[oldName]
	delete(c.variables, oldName)

	// Update the event triggers
	for _, event := range events {
		c.listeners[event][newName] = c.listeners[event][oldName]
		delete(c.listeners[event], oldName)
	}

	// Trigger the event
	c.eventTriggered(RenameVarEvent, newName)
}

func (c *Controller) Delete(name string) {
	// Check if the variable exists
	if _, exists := c.variables[name]; !exists {
		log.Println("Error: attempt to delete a variable that doesn't exist:", name)
		return
	}

	// Update the variable map
	delete(c.variables, name)

	// Trigger the event
	c.eventTriggered(DeleteVarEvent, name)

	// Delete the variable from the listener map
	for _, event := range events {
		delete(c.listeners[event], name)
	}
}

func (c *Controller) AddListener(event Event, variableName string, callback func(string)) {
	if _, exists := c.listeners[event][variableName]; !exists {
		c.listeners[event][variableName] = make([]func(string), 0)
	}
	c.listeners[event][variableName] = append(c.listeners[event][variableName], callback)
}

//func (c *Controller) DeleteListener(event Event, variableName string, callback func(string)) {
//callbacks, exists := c.listeners[event][variableName]
//if !exists {
//return
//}
//newCallbacks := make([]func(string), 0, len(callbacks-1))
//callbackFound := false
//for index, cb := range callbacks {
//if callback == cb {
//callbackFound = true
//}
//}

//c.listeners[event][variableName] = append(c.listeners[event][variableName], callback)
//}

func (c Controller) eventTriggered(event Event, variableName string) {
	callbacks, exists := c.listeners[event][variableName]
	if exists {
		for _, callback := range callbacks {
			callback(variableName)
		}
	}
	for _, callback := range c.listeners[event]["*"] {
		callback(variableName)
	}
}
