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
		objects:         make(map[*ObjectHolder]bool),
		objectNames:     make(map[string]*ObjectHolder),
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

func (c *Controller) NewObject(name string, obj *Object) *ObjectHolder {
	// Create the object holder
	holder := &ObjectHolder{object: obj, name: name, controller: c}
	_, isProducer := (*obj).(Producer)
	if isProducer {
		holder.dependencies = make([]*ObjectHolder, 0)
	}
	_, isConsumer := (*obj).(Consumer)
	if isConsumer {
		holder.dependents = make([]*ObjectHolder, 0)
	}

	// Add the object to the map
	c.objects[holder] = true
	c.objectNames[name] = holder

	// Trigger the callback
	for callback := range c.globalListeners[NewObjectEvent] {
		(*callback)(holder)
	}
	return holder
}

func (c *Controller) RemoveObject(holder *ObjectHolder) {
	(*holder.object).Close()

	// Remove object from the maps
	if _, exists := c.objectNames[holder.Name()]; exists {
		delete(c.objectNames, holder.Name())
	}
	if _, exists := c.objects[holder]; exists {
		delete(c.objects, holder)
	}

	// Run the event
	c.EventTriggered(DeleteObjectEvent, holder)
	for _, event := range events {
		delete(c.listeners[event], holder)
	}
}

func (c *Controller) AddListener(event Event, holder *ObjectHolder, callback *func(*ObjectHolder)) {
	if _, exists := c.listeners[event]; !exists {
		// Should I just add the event?
		return
	}
	if _, exists := c.listeners[event][holder]; !exists {
		c.listeners[event][holder] = make(map[Listener]bool)
	}
	if !c.listeners[event][holder][callback] {
		c.listeners[event][holder][callback] = true
	}
}

func (c *Controller) AddGlobalListener(event Event, callback *func(*ObjectHolder)) {
	if _, exists := c.globalListeners[event]; !exists {
		// Should I just add the event?
		return
	}
	if !c.globalListeners[event][callback] {
		c.globalListeners[event][callback] = true
	}
}

func (c *Controller) DeleteListener(event Event, holder *ObjectHolder, callback *func(*ObjectHolder)) {
	if _, exists := c.listeners[event][holder]; !exists {
		return
	}
	if _, exists := c.listeners[event][holder][callback]; !exists {
		return
	}
	delete(c.listeners[event][holder], callback)
}

func (c Controller) EventTriggered(event Event, holder *ObjectHolder) {
	if _, exists := c.listeners[event]; !exists {
		log.Println("Error: event", event, "does not exist!")
		return
	}
	callbacks, exists := c.listeners[event][holder]
	if exists {
		for callback := range callbacks {
			(*callback)(holder)
		}
	}
	for callback := range c.globalListeners[event] {
		(*callback)(holder)
	}
}
