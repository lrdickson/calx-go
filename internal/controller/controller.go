package controller

import (
	"fmt"
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
	objectIdCount   ObjectId
	listeners       listenerMap
	globalListeners map[Event]map[Listener]bool
}

func NewController() *Controller {
	// Make the new controller
	controller := &Controller{
		objects:         make(map[ObjectId]*Object),
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

func (c *Controller) NewObject(name string) ObjectId {
	// Create the object object
	objectId := c.objectIdCount
	c.objectIdCount++
	object := &Object{name: name, id: objectId}

	// Add the object to the map
	c.objects[objectId] = object
	c.objectNames[name] = objectId

	// Trigger the callback
	for callback := range c.globalListeners[NewObjectEvent] {
		(*callback)(object)
	}
	return objectId
}

func (c *Controller) RemoveObject(id ObjectId) {
	// Remove object from the maps
	object, exists := c.objects[id]
	if exists {
		delete(c.objects, id)
	}
	if _, exists := c.objectNames[object.Name()]; exists {
		delete(c.objectNames, object.Name())
	}

	// Run the event
	c.EventTriggered(DeleteObjectEvent, id)
	for _, event := range events {
		delete(c.listeners[event], id)
	}
}

func (c *Controller) Rename(id ObjectId, name string) error {
	// Make sure the object exists
	o, exists := c.objects[id]
	if !exists {
		return fmt.Errorf("Name object not found with id: %d", id)
	}

	// Make sure the name is unique
	if _, exists := c.objectNames[name]; exists {
		return fmt.Errorf("The name %s is taken", name)
	}

	// Make sure the name is valid
	if err := NameValid(name); err != nil {
		return err
	}

	// Update the name
	if _, exists := c.objectNames[o.name]; exists {
		delete(c.objectNames, o.name)
	}
	c.objectNames[name] = c.objectNames[o.name]
	o.name = name
	c.EventTriggered(RenameObjectEvent, o.id)

	// Success
	return nil
}

func (c *Controller) Objects(id ObjectId) *Object {
	return c.objects[id]
}

func (c Controller) IterObjects(iter func(ObjectId, *Object) bool) {
	for id, obj := range c.objects {
		cont := iter(id, obj)
		if !cont {
			break
		}
	}
}

func (c Controller) ObjectCount() int {
	return len(c.objects)
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

func (c *Controller) AddListener(event Event, id ObjectId, callback *func(*Object)) {
	if _, exists := c.listeners[event]; !exists {
		// Should I just add the event?
		return
	}
	if _, exists := c.listeners[event][id]; !exists {
		c.listeners[event][id] = make(map[Listener]bool)
	}
	if !c.listeners[event][id][callback] {
		c.listeners[event][id][callback] = true
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

func (c *Controller) DeleteListener(event Event, id ObjectId, callback *func(*Object)) {
	if _, exists := c.listeners[event][id]; !exists {
		return
	}
	if _, exists := c.listeners[event][id][callback]; !exists {
		return
	}
	delete(c.listeners[event][id], callback)
}

func (c Controller) EventTriggered(event Event, id ObjectId) {
	if _, exists := c.listeners[event]; !exists {
		log.Println("Error: event", event, "does not exist!")
		return
	}
	object, exists := c.objects[id]
	if !exists {
		log.Println("Error: object", event, "does not exist!")
		return
	}
	callbacks, exists := c.listeners[event][id]
	if exists {
		for callback := range callbacks {
			(*callback)(object)
		}
	}
	for callback := range c.globalListeners[event] {
		(*callback)(object)
	}
}
