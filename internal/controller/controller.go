package controller

import (
	"log"
	"strconv"

	"github.com/fxamacker/cbor/v2"
)

type Event int

const (
	NewObjectEvent Event = iota + 1
	RenameObjectEvent
	DeleteObjectEvent
)

var events []Event = []Event{NewObjectEvent, RenameObjectEvent, DeleteObjectEvent}

type Listener *func(ObjectId)
type listenerMap map[Event]map[ObjectId]map[Listener]bool

type Controller struct {
	objectIdCount     ObjectId
	objects           map[ObjectId]*Object
	objectNames       map[string]ObjectId
	objectOutput      map[OutputVersion]map[ObjectId]cbor.RawMessage
	listeners         listenerMap
	globalListeners   map[Event]map[Listener]bool
	metadataListeners map[ObjectId]map[string]map[Listener]bool
	latestVersion     OutputVersion
	completeVersion   OutputVersion
}

func NewController() *Controller {
	// Make the new controller
	controller := &Controller{
		objects:           make(map[ObjectId]*Object),
		objectNames:       make(map[string]ObjectId),
		listeners:         make(listenerMap),
		globalListeners:   make(map[Event]map[Listener]bool),
		metadataListeners: make(map[ObjectId]map[string]map[Listener]bool),
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

func (c *Controller) AddListener(event Event, id ObjectId, callback Listener) {
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

func (c *Controller) DeleteListener(event Event, id ObjectId, callback Listener) {
	if _, exists := c.listeners[event][id]; !exists {
		return
	}
	if _, exists := c.listeners[event][id][callback]; !exists {
		return
	}
	delete(c.listeners[event][id], callback)
}

func (c *Controller) AddGlobalListener(event Event, callback Listener) {
	if _, exists := c.globalListeners[event]; !exists {
		// Should I just add the event?
		return
	}
	if !c.globalListeners[event][callback] {
		c.globalListeners[event][callback] = true
	}
}

func (c *Controller) DeleteGlobalListener(event Event, callback Listener) {
	if _, exists := c.globalListeners[event][callback]; !exists {
		return
	}
	delete(c.globalListeners[event], callback)
}

func (c Controller) EventTriggered(event Event, id ObjectId) {
	if _, exists := c.listeners[event]; !exists {
		log.Println("Error: event", event, "does not exist!")
		return
	}
	if _, exists := c.objects[id]; !exists {
		log.Println("Error: object", event, "does not exist!")
		return
	}
	callbacks, exists := c.listeners[event][id]
	if exists {
		for callback := range callbacks {
			(*callback)(id)
		}
	}
	for callback := range c.globalListeners[event] {
		(*callback)(id)
	}
}
