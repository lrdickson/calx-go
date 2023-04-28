package controller

import (
	"encoding/json"
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

type ObjectId int

type Object struct {
	id           ObjectId
	name         string
	metadata     map[string]json.RawMessage
	dependencies []*Object
	dependents   []*Object
}

func (c *Controller) NewObject(name string) ObjectId {
	// Create the object
	objectId := c.objectIdCount
	c.objectIdCount++
	obj := &Object{
		id:           objectId,
		name:         name,
		metadata:     make(map[string]json.RawMessage),
		dependencies: make([]*Object, 0),
		dependents:   make([]*Object, 0),
	}

	// Add the object to the map
	c.objects[objectId] = obj
	c.objectNames[name] = objectId

	// Trigger the callback
	for callback := range c.globalListeners[NewObjectEvent] {
		(*callback)(objectId)
	}
	return objectId
}

func (c *Controller) RemoveObject(id ObjectId) {
	// Remove obj from the maps
	obj, exists := c.objects[id]
	if exists {
		delete(c.objects, id)
	}
	if _, exists := c.objectNames[obj.name]; exists {
		delete(c.objectNames, obj.name)
	}

	// Run the event
	c.EventTriggered(DeleteObjectEvent, id)
	for _, event := range events {
		delete(c.listeners[event], id)
	}
}

func (c *Controller) getObject(id ObjectId) (*Object, error) {
	obj, exists := c.objects[id]
	if !exists {
		return obj, errors.New("Object does not exist")
	}
	return obj, nil
}

func (c *Controller) Name(id ObjectId) (string, error) {
	obj, exists := c.objects[id]
	if !exists {
		return "", errors.New("Object does not exist")
	}
	return obj.name, nil
}

func (c *Controller) SetName(id ObjectId, name string) error {
	// Make sure the object exists
	o, err := c.getObject(id)
	if err != nil {
		return err
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

func (c *Controller) MetaData(id ObjectId, key string) (json.RawMessage, error) {
	obj, err := c.getObject(id)
	if err != nil {
		return json.RawMessage{}, err
	}
	data, exists := obj.metadata[key]
	if !exists {
		return json.RawMessage{}, fmt.Errorf("Object %d(%s) does not have key %s",
			id, obj.name, key)
	}
	return data, nil
}

func (c *Controller) SetMetaData(id ObjectId, key string, data json.RawMessage) error {
	_, err := c.getObject(id)
	if err != nil {
		return err
	}
	c.objects[id].metadata[key] = data
	return nil
}

func (c *Controller) Output(id ObjectId) (json.RawMessage, error) {
	obj, err := c.getObject(id)
	if err != nil {
		return json.RawMessage{}, err
	}
	return obj.output, nil
}

func (c *Controller) SetOutput(id ObjectId, data json.RawMessage, version OutputVersion) error {
	_, err := c.getObject(id)
	if err != nil {
		return err
	}
	c.objects[id].output[version] = data
	return nil
}

func (c Controller) IterObjects(iter func(ObjectId) bool) {
	for id := range c.objects {
		cont := iter(id)
		if !cont {
			break
		}
	}
}

func (c Controller) ObjectCount() int {
	return len(c.objects)
}
