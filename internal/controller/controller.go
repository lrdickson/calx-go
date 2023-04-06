package controller

type Event int

const (
	New Event = iota
	Rename
	Delete
)

type Controller struct {
	variables map[string]*Variable
}

func NewController() *Controller {
	return &Controller{
		variables: make(map[string]*Variable),
	}
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