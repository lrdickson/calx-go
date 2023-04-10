package controller

import (
	"log"
	"strconv"

	"github.com/lrdickson/calx/internal/variable"
)

type Event string

const (
	NewVarEvent    Event = "NewVar"
	RenameVarEvent Event = "RenameVar"
	DeleteVarEvent Event = "DeleteVar"
)

var events []Event = []Event{NewVarEvent, RenameVarEvent, DeleteVarEvent}

type ListenerId int64

// listeners[event][variableName][listenerId]
type listenerMap map[Event]map[string]map[ListenerId]func(string)

type Controller struct {
	variables     map[string]*variable.Variable
	variableCount uint64
	listeners     listenerMap
	listenerCount ListenerId
}

func NewController() *Controller {
	// Make the new controller
	controller := &Controller{
		variables:     make(map[string]*variable.Variable),
		variableCount: 1,
		listeners:     make(listenerMap),
		listenerCount: 1,
	}

	// Initialize the listeners map
	for _, event := range events {
		controller.AddEvent(event)
	}

	// Create the controller
	return controller
}

func (c *Controller) AddEvent(event Event) {
	if _, exists := c.listeners[event]; exists {
		log.Println("Error: event", event, "already exists!")
		return
	}
	c.listeners[event] = make(map[string]map[ListenerId]func(string))
	// Add universal listenner
	c.listeners[event]["*"] = make(map[ListenerId]func(string))
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

func (c *Controller) UniqueName() string {
	name := ""
	for {
		name = "var" + strconv.FormatUint(c.variableCount, 10)
		if _, taken := c.variables[name]; taken {
			c.variableCount++
		} else {
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
	c.AddVariable(c.UniqueName(), &formula)
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
	c.EventTriggered(RenameVarEvent, newName)
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
	c.EventTriggered(DeleteVarEvent, name)

	// Delete the variable from the listener map
	for _, event := range events {
		delete(c.listeners[event], name)
	}
}

func (c *Controller) AddListener(event Event, variableName string, callback func(string)) ListenerId {
	if _, exists := c.listeners[event][variableName]; !exists {
		c.listeners[event][variableName] = make(map[ListenerId]func(string))
	}
	listenerId := c.listenerCount
	c.listeners[event][variableName][listenerId] = callback
	c.listenerCount++
	return listenerId
}

func (c *Controller) DeleteListener(event Event, variableName string, listenerId ListenerId) {
	if _, exists := c.listeners[event][variableName]; !exists {
		return
	}
	if _, exists := c.listeners[event][variableName][listenerId]; !exists {
		return
	}
	delete(c.listeners[event][variableName], listenerId)
}

func (c Controller) EventTriggered(event Event, variableName string) {
	if _, exists := c.listeners[event]; !exists {
		log.Println("Error: event", event, "does not exist!")
		return
	}
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
