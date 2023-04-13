package controller

import (
	"log"
	"strconv"
)

type Event string

const (
	NewVarEvent    Event = "NewVar"
	RenameVarEvent Event = "RenameVar"
	DeleteVarEvent Event = "DeleteVar"
)

var events []Event = []Event{NewVarEvent, RenameVarEvent, DeleteVarEvent}

type ListenerId int64
type ObjectId int64

// listeners[event][variableName][listenerId]
type listenerMap map[Event]map[ObjectId]map[ListenerId]func(string)

type Controller struct {
	objects         map[ObjectId]*Object
	objectIdCount   ObjectId
	objectNames     map[string]ObjectId
	objectNameCount uint64
	listeners       listenerMap
	listenerIdCount ListenerId
}

func NewController() *Controller {
	// Make the new controller
	controller := &Controller{
		objects:         make(map[ObjectId]*Object),
		objectIdCount:   1,
		objectNames:     make(map[string]ObjectId),
		objectNameCount: 1,
		listeners:       make(listenerMap),
		listenerIdCount: 1,
	}

	// Initialize the listeners map
	for _, event := range events {
		controller.AddEvent(event)
	}

	// Create the controller
	return controller
}

func (c *Controller) AddEvent(event Event) {
	// Add the event
	if _, exists := c.listeners[event]; exists {
		log.Println("Error: event", event, "already exists!")
		return
	}
	c.listeners[event] = make(map[ObjectId]map[ListenerId]func(string))

	// Add universal listener
	universalName := "*"
	c.objectNames[universalName] = c.objectIdCount
	c.listeners[event][c.objectIdCount] = make(map[ListenerId]func(string))
	c.objectIdCount++
}

func (c Controller) IterVariables(iter func(ObjectId, *Object) bool) {
	for key, value := range c.objects {
		cont := iter(key, value)
		if !cont {
			break
		}
	}
}

func (c Controller) Variables(id ObjectId) *Object {
	return c.objects[id]
}

func (c Controller) VariableCount() int {
	return len(c.objects)
}

func (c *Controller) UniqueName() string {
	name := ""
	for {
		name = "obj" + strconv.FormatUint(c.objectNameCount, 10)
		if _, taken := c.objectNames[name]; taken {
			c.objectIdCount++
		} else {
			break
		}
	}
	return name
}

func (c *Controller) AddObject(name string, obj *Object) {
	// Add the object to the map
	c.objects[c.objectIdCount] = obj
	c.objectNames[name] = c.objectIdCount
	c.objectIdCount++

	// Trigger the callback
	for _, callback := range c.listeners[NewVarEvent][c.objectNames["*"]] {
		callback(name)
	}
}

func (c *Controller) Rename(oldName, newName string) {
	// Check if the oldName exists
	if _, exists := c.objectNames[oldName]; !exists {
		log.Println("Error: attempt to rename a variable that doesn't exist:", oldName)
		return
	}

	// Update the variable map
	c.objectNames[newName] = c.objectNames[oldName]
	delete(c.objectNames, oldName)

	// Trigger the event
	c.EventTriggered(RenameVarEvent, newName)
}

func (c *Controller) Delete(name string) {
	// Check if the variable exists
	if _, exists := c.objects[name]; !exists {
		log.Println("Error: attempt to delete a variable that doesn't exist:", name)
		return
	}

	// Delete the variable
	(*c.objects[name]).Close()
	delete(c.objects, name)

	// Run the event
	c.EventTriggered(DeleteVarEvent, name)
	for _, event := range events {
		delete(c.listeners[event], name)
	}
}

func (c *Controller) AddListener(event Event, variableName string, callback func(string)) ListenerId {
	if _, exists := c.listeners[event][variableName]; !exists {
		c.listeners[event][variableName] = make(map[ListenerId]func(string))
	}
	listenerId := c.listenerIdCount
	c.listeners[event][variableName][listenerId] = callback
	c.listenerIdCount++
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
