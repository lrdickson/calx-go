package controller

import "testing"

func AddBaseObject(c *Controller) *ObjectHolder {
	var obj Object = &BaseObject{}
	return c.NewObject(c.UniqueName(), &obj)
}

func TestAddDeleteListener(t *testing.T) {
	// Test successful listener add
	c := NewController()
	obj := AddBaseObject(c)
	callback := func(_ *ObjectHolder) {}
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
	callback := func(obj *ObjectHolder) {
		name = obj.Name()
	}
	c.AddGlobalListener(NewObjectEvent, &callback)
	obj := AddBaseObject(c)

	// Check for a successful add
	AddBaseObject(c)
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
	obj := AddBaseObject(c)

	// Add the listener
	listenerCalled := false
	callback := func(_ *ObjectHolder) {
		listenerCalled = true
	}
	c.AddListener(RenameObjectEvent, obj, &callback)

	// Rename the variable
	newName := "NewName"
	obj.SetName(newName)

	// Check the result
	if obj.Name() != newName {
		t.Fatal("Rename failed")
	}
	if !listenerCalled {
		t.Fatal("Rename callback not called")
	}
}
