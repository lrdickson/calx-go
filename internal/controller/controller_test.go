package controller

import "testing"

func TestAddDeleteListener(t *testing.T) {
	// Test successful listener add
	c := NewController()
	varName := "*"
	listenerId := c.AddListener(NewVarEvent, varName, func(variableName string) {})
	if _, exists := c.listeners[NewVarEvent][varName][listenerId]; !exists {
		t.Fatal("Failed to add a listener")
	}

	// Test successful listener delete
	c.DeleteListener(NewVarEvent, varName, listenerId)
	if _, exists := c.listeners[NewVarEvent][varName][listenerId]; exists {
		t.Fatal("Failed to delete a listener")
	}
}

func TestAddDeleteVar(t *testing.T) {
	// Setup the add listener
	c := NewController()
	name := ""
	c.AddListener(NewVarEvent, "*", func(variableName string) {
		name = variableName
	})

	// Check for a successful add
	c.AddFormula()
	if name == "" {
		t.Fatal("Name reciever is still empty")
	}
	if _, exists := c.variables[name]; !exists {
		t.Fatal(name, "not in variables")
	}

	// Check for a successful delete
	c.Delete(name)
	if _, exists := c.variables[name]; exists {
		t.Fatal(name, "still in variables after delete")
	}
}

func TestRenameVar(t *testing.T) {
	// Add a variable
	c := NewController()
	name := ""
	c.AddListener(NewVarEvent, "*", func(variableName string) {
		name = variableName
	})
	c.AddFormula()

	// Rename the variable
	newName := "NewName"
	c.Rename(name, newName)

	// Check the result
	if _, exists := c.variables[newName]; !exists {
		t.Fatal(newName, "not in variables")
	}
}
