package controller

import (
	"log"
	"strconv"
)

type Event string

const (
	NewObjectEvent    Event = "NewObject"
	RenameObjectEvent Event = "RenameObject"
	DeleteObjectEvent Event = "DeleteObject"
)

var events []Event = []Event{NewObjectEvent, RenameObjectEvent, DeleteObjectEvent}

type Listener *func(*Object)
type listenerMap map[Event]map[ObjectId]map[Listener]bool

type Controller struct {
	objects         map[ObjectId]*Object
	objectNames     map[string]ObjectId
	listeners       listenerMap
	globalListeners map[Event]map[Listener]bool
}

func NewController() *Controller {
	// Make the new controller
	controller := &Controller{
		objects:         make(map[ObjectId]bool),
		objectNames:     make(map[string]ObjectId),
		listeners:       make(listenerMap),
		globalListeners: make(map[Event]map[Listener]bool),
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
	c.listeners[event] = make(map[ObjectId]map[Listener]bool)
	c.globalListeners[event] = make(map[Listener]bool)
}

func (c Controller) IterVariables(iter func(ObjectId) bool) {
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

func (c *Controller) NewObject(name string) *Object {
	// Create the object object
	object := &Object{name: name, controller: c}

	// Add the object to the map
	c.objects[object] = true
	c.objectNames[name] = object

	// Trigger the callback
	for callback := range c.globalListeners[NewObjectEvent] {
		(*callback)(object)
	}
	return object
}

func (c *Controller) RemoveObject(object *Object) {
	(*object.engine).Close()

	// Remove object from the maps
	if _, exists := c.objectNames[object.Name()]; exists {
		delete(c.objectNames, object.Name())
	}
	if _, exists := c.objects[object]; exists {
		delete(c.objects, object)
	}

	// Run the event
	c.EventTriggered(DeleteObjectEvent, object)
	for _, event := range events {
		delete(c.listeners[event], object)
	}
}

func (c *Controller) AddListener(event Event, object *Object, callback *func(*Object)) {
	if _, exists := c.listeners[event]; !exists {
		// Should I just add the event?
		return
	}
	if _, exists := c.listeners[event][object]; !exists {
		c.listeners[event][object] = make(map[Listener]bool)
	}
	if !c.listeners[event][object][callback] {
		c.listeners[event][object][callback] = true
	}
}

func (c *Controller) AddGlobalListener(event Event, callback *func(*Object)) {
	if _, exists := c.globalListeners[event]; !exists {
		// Should I just add the event?
		return
	}
	if !c.globalListeners[event][callback] {
		c.globalListeners[event][callback] = true
	}
}

func (c *Controller) DeleteListener(event Event, object *Object, callback *func(*Object)) {
	if _, exists := c.listeners[event][object]; !exists {
		return
	}
	if _, exists := c.listeners[event][object][callback]; !exists {
		return
	}
	delete(c.listeners[event][object], callback)
}

func (c Controller) EventTriggered(event Event, object *Object) {
	if _, exists := c.listeners[event]; !exists {
		log.Println("Error: event", event, "does not exist!")
		return
	}
	callbacks, exists := c.listeners[event][object]
	if exists {
		for callback := range callbacks {
			(*callback)(object)
		}
	}
	for callback := range c.globalListeners[event] {
		(*callback)(object)
	}
}
