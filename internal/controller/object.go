package controller

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/fxamacker/cbor/v2"
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

type ObjectId int64
type ObjectVersion int64

type Object struct {
	id           ObjectId
	name         string
	metadata     map[string]cbor.RawMessage
	dependencies []ObjectId
	dependents   []ObjectId
	output       cbor.RawMessage
}

func (c *Controller) NewObject(name string) ObjectId {
	// Create the object
	objectId := c.objectIdCount
	c.objectIdCount++
	obj := &Object{
		id:           objectId,
		name:         name,
		metadata:     make(map[string]cbor.RawMessage),
		dependencies: make([]ObjectId, 0),
		dependents:   make([]ObjectId, 0),
	}

	// Add the object to the map
	c.objects[c.latestVersion][objectId] = obj
	c.objectNames[name] = objectId

	// Trigger the callback
	for callback := range c.globalListeners[NewObjectEvent] {
		(*callback)(objectId)
	}
	return objectId
}

func (c *Controller) RemoveObject(id ObjectId) error {
	// Initialize the new generation
	previousVersion := c.latestVersion
	c.latestVersion++
	newGeneration := make(map[ObjectId]*Object)

	descendants, err := c.descendants(previousVersion, id)
	if err != nil {
		return err
	}
	for id, obj := range c.objects[previousVersion] {
		if descendants[id] {
			newObj := &Object{
				id:           obj.id,
				name:         obj.name,
				dependencies: obj.dependencies,
				dependents:   obj.dependents,
			}
			newGeneration[id] = newObj
		} else {
			newGeneration[id] = obj
		}
	}
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

func (c *Controller) getObject(version ObjectVersion, id ObjectId) (*Object, error) {
	_, exists := c.objects[version]
	if !exists {
		return nil, fmt.Errorf("Version %d not currently available", version)
	}

	obj, exists := c.objects[version][id]
	if !exists {
		return obj, fmt.Errorf("Object %d does not exist in version %d", id, version)
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

func (c *Controller) MetaData(id ObjectId, key string) (cbor.RawMessage, error) {
	obj, err := c.getObject(id)
	if err != nil {
		return cbor.RawMessage{}, err
	}
	data, exists := obj.metadata[key]
	if !exists {
		return cbor.RawMessage{}, fmt.Errorf("Object %d(%s) does not have key %s",
			id, obj.name, key)
	}
	return data, nil
}

func (c *Controller) SetMetaData(id ObjectId, key string, data cbor.RawMessage) error {
	_, err := c.getObject(id)
	if err != nil {
		return err
	}
	c.objects[id].metadata[key] = data
	return nil
}

func (c *Controller) getOutputMap(version ObjectVersion) (map[ObjectId]cbor.RawMessage, error) {
	outputMap, exists := c.objectOutput[version]
	if !exists {
		return outputMap, fmt.Errorf("Version %d is not currently available", version)
	}
	return outputMap, nil
}

func (c *Controller) Output(version ObjectVersion, id ObjectId) (cbor.RawMessage, error) {
	outputMap, err := c.getOutputMap(version)
	if err != nil {
		return cbor.RawMessage{}, err
	}
	output, exists := outputMap[id]
	if !exists {
		return output, fmt.Errorf("Version %d does not contain output for object %d", version, id)
	}
	return output, nil
}

func (c *Controller) SetOutput(version ObjectVersion, id ObjectId, data cbor.RawMessage) error {
	outputMap, err := c.getOutputMap(version)
	if err != nil {
		return err
	}
	_, exists := outputMap[id]
	if !exists {
		return fmt.Errorf("Version %d does not contain output for object %d", version, id)
	}
	outputMap[id] = data
	return nil
}

func (c Controller) Range(iter func(ObjectId) bool) {
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

func (c *Controller) descendants(version ObjectVersion, id ObjectId) (map[ObjectId]bool, error) {
	obj, err := c.getObject(version, id)
	if err != nil {
		return nil, err
	}
	descendants := map[ObjectId]bool{}
	for _, dependent := range obj.dependents {
		descendants[dependent] = true
		grandDescendants, err := c.descendants(version, dependent)
		if err != nil {
			return nil, err
		}
		for descendant := range grandDescendants {
			descendants[descendant] = true
		}
	}
	return descendants, nil
}

func (c *Controller) newGeneration(changedObj ObjectId) error {
	// Initialize the new generation
	previousVersion := c.latestVersion
	c.latestVersion++
	newGeneration := make(map[ObjectId]*Object)

	descendants, err := c.descendants(previousVersion, changedObj)
	if err != nil {
		return err
	}
	for id, obj := range c.objects[previousVersion] {
		if descendants[id] {
			newObj := &Object{
				id:           obj.id,
				name:         obj.name,
				dependencies: obj.dependencies,
				dependents:   obj.dependents,
			}
			newGeneration[id] = newObj
		} else {
			newGeneration[id] = obj
		}
	}
	c.objects[c.latestVersion] = newGeneration
	return nil
}
