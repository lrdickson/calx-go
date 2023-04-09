package controller

import "testing"

func TestAddFormula(t *testing.T) {
	// Setup the add listener
	c := NewController()
	nameReceiver := ""
	c.AddListener(NewVarEvent, "*", func(variableName string) {
		nameReceiver = variableName
	})

	// Add a variable
	c.AddFormula()

	// Check the output
	if nameReceiver == "" {
		t.Fatal("Name reciever is still empty")
	}
}

func TestRename(t *testing.T) {
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
