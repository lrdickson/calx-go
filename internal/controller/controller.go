package controller

import (
	"errors"
	"fmt"
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

type ObjectHolder struct {
	object       *Object
	dependencies map[*ObjectHolder]bool
	dependents   map[*ObjectHolder]bool
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

	// Record the name
	o.name = name

	// Success
	return nil
}

type Listener *func(*ObjectHolder)
type listenerMap map[Event]map[*ObjectHolder]map[Listener]bool

type Controller struct {
	objects         map[*ObjectHolder]bool
	objectNames     map[string]*ObjectHolder
	listeners       listenerMap
	globalListeners map[Event]map[Listener]bool
}

func NewController() *Controller {
	// Make the new controller
	controller := &Controller{
		objects:     make(map[*ObjectHolder]bool),
		objectNames: make(map[string]*ObjectHolder),
		listeners:   make(listenerMap),
	}

	// Initialize the listeners map
	for _, event := range events {
		controller.AddEvent(event)
	}

	// Create the controller
	return controller
}

func (c *Controller) AddEvent(event Event) {
	// Prevent duplicate listener maps
	if _, exists := c.listeners[event]; exists {
		log.Println("Error: event", event, "already exists!")
		return
	}

	// Add the event to the listener maps
	c.listeners[event] = make(map[*ObjectHolder]map[Listener]bool)
	c.globalListeners[event] = make(map[Listener]bool)
}

func (c Controller) IterVariables(iter func(*ObjectHolder) bool) {
	for key := range c.objects {
		cont := iter(key)
		if !cont {
			break
		}
	}
}

func (c Controller) VariableCount() int {
	return len(c.objects)
}

func (c *Controller) UniqueName() string {
	name := ""
	objectNameCount := 1
	for {
		name = "obj" + strconv.Itoa(objectNameCount)
		if _, taken := c.objectNames[name]; taken {
			objectNameCount++
		} else {
			break
		}
	}
	return name
}

func (c *Controller) AddObject(name string, obj *Object) error {
	// Make sure the object can do something
	_, isProducer := (*obj).(Producer)
	_, isConsumer := (*obj).(Consumer)
	if !isProducer && !isConsumer {
		return errors.New("Provided object is neither a producer nor a consumer")
	}

	// Add the object to the map
	c.objects[c.objectIdCount] = obj
	c.objectNames[name] = c.objectIdCount
	c.objectIdCount++

	// Trigger the callback
	for _, callback := range c.listeners[NewVarEvent][c.objectNames["*"]] {
		callback(name)
	}
	return nil
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
