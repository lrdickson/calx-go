package controller

import (
	"testing"
)

func AddObject(c *Controller) ObjectId {
	return c.NewObject(c.UniqueName())
}

func TestAddDeleteListener(t *testing.T) {
	// Test successful listener add
	c := NewController()
	obj := AddObject(c)
	callback := func(_ *Object) {}
	c.AddListener(NewObjectEvent, obj, &callback)
	if _, exists := c.listeners[NewObjectEvent][obj][&callback]; !exists {
		t.Fatal("Failed to add a listener")
	}

	// Test successful listener delete
	c.DeleteListener(NewObjectEvent, obj, &callback)
	if _, exists := c.listeners[NewObjectEvent][obj][&callback]; exists {
		t.Fatal("Failed to delete a listener")
	}
}

func TestAddDeleteVar(t *testing.T) {
	// Setup the add listener
	c := NewController()
	name := ""
	callback := func(obj *Object) {
		name = obj.Name()
	}
	c.AddGlobalListener(NewObjectEvent, &callback)
	obj := AddObject(c)

	// Check for a successful add
	AddObject(c)
	if name == "" {
		t.Fatal("Name reciever is still empty")
	}
	if _, exists := c.objects[obj]; !exists {
		t.Fatal("Object not added to controller object map")
	}

	// Check for a successful delete
	c.RemoveObject(obj)
	if _, exists := c.objects[obj]; exists {
		t.Fatal(name, "still in variables after remove")
	}
}

func TestRenameVar(t *testing.T) {
	// Add a variable
	c := NewController()
	var obj *Object
	newObjectCallback := func(o *Object) {
		obj = o
	}
	c.AddGlobalListener(NewObjectEvent, &newObjectCallback)
	id := AddObject(c)

	// Add the listener
	listenerCalled := false
	callback := func(_ *Object) {
		listenerCalled = true
	}
	c.AddListener(RenameObjectEvent, id, &callback)

	// Rename the variable
	newName := "NewName"
	err := c.Rename(id, newName)
	if err != nil {
		t.Fatal("Failed to rename variable:", err)
	}

	// Check the result
	if obj.Name() != newName {
		t.Fatal("Rename failed")
	}
	if !listenerCalled {
		t.Fatal("Rename callback not called")
	}
}
