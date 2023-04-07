package controller

import (
	"strconv"

	"github.com/lrdickson/calx/internal/variable"
)

type Event int

const (
	New Event = iota
	Rename
	Delete
)

type Controller struct {
	variables     map[string]*variable.Variable
	variableCount int
}

func NewController() *Controller {
	return &Controller{
		variables: make(map[string]*variable.Variable),
	}
}

func (c Controller) IterVariables(iter func(string, *variable.Variable) bool) {
	for key, value := range c.variables {
		cont := iter(key, value)
		if !cont {
			break
		}
	}
}

func (c Controller) Variables(name string) *variable.Variable {
	return c.variables[name]
}

func (c *Controller) uniqueName() string {
	name := ""
	for {
		name = "var" + strconv.Itoa(c.variableCount)
		c.variableCount++
		if _, taken := c.variables[name]; !taken {
			break
		}
	}
	return name
}

func (c *Controller) AddFormula() {
	name := c.uniqueName()
	var formula variable.Variable = variable.Formula{}
	c.variables[name] = &formula
}

func (c *Controller) Rename(oldName, newName string) {
	// Check if the oldName exists
	if _, exists := c.variables[oldName]; !exists {
		return
	}

	// Update the variable map
	c.variables[newName] = c.variables[oldName]
	delete(c.variables, oldName)
}
