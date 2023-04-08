package controller

import (
	"log"
	"strconv"

	"github.com/lrdickson/calx/internal/variable"
)

type Event int

const (
	NewEvent Event = iota
	RenameEvent
	DeleteEvent
)

var events []Event = []Event{NewEvent, RenameEvent, DeleteEvent}

// listeners[event][variableName]
type listenerMap map[Event]map[string][]func()

type Controller struct {
	variables     map[string]*variable.Variable
	variableCount int
	listeners     listenerMap
}

func NewController() *Controller {
	// Initialize the listeners map
	listeners := make(listenerMap)
	for _, event := range events {
		listeners[event] = make(map[string][]func())
		// Add universal listenner
		listeners[event]["*"] = make([]func(), 0)
	}

	// Create the controller
	return &Controller{
		variables: make(map[string]*variable.Variable),
		listeners: listeners,
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

func (c *Controller) AddFormula() {
	name := c.uniqueName()
	var formula variable.Variable = variable.Formula{}
	c.variables[name] = &formula
	c.eventTriggered(NewEvent, "*")
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

	// Trigger the event
	c.eventTriggered(RenameEvent, oldName)
	c.eventTriggered(RenameEvent, "*")
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
	c.eventTriggered(RenameEvent, name)
	c.eventTriggered(RenameEvent, "*")

	// Delete the variable from the listener map
	for _, event := range events {
		delete(c.listeners[event], name)
	}
}

func (c *Controller) AddListener(event Event, variableName string, callback func()) {
	if _, exists := c.listeners[event][variableName]; !exists {
		c.listeners[event][variableName] = make([]func(), 0)
	}
	c.listeners[event][variableName] = append(c.listeners[event][variableName], callback)
}

func (c Controller) eventTriggered(event Event, variableName string) {
	callbacks, exists := c.listeners[event][variableName]
	if !exists {
		log.Println("Error: attempt to trigger an event for a variable that doesn't exist:", variableName)
		return
	}
	for _, callback := range callbacks {
		callback()
	}
}
