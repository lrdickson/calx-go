package formula

import "github.com/lrdickson/calx/internal/controller"

type Formula struct {
	controller.BaseObject
	code         string
	dependencies map[string]bool
	dependents   map[string]bool
}

func (f *Formula) Code() string {
	return f.code
}

func AddFormula(c *controller.Controller) {
	var formula controller.Object = &Formula{}
	controller.AddObject(c, c.UniqueName(), &formula)
}

func (f *Formula) AddDependency(c *controller.Controller, name string) {
	f.dependencies[name] = true
	c.AddListener(controller.RenameVarEvent, name, func(dependencyName string) {
		f.dependencies[dependencyName] = true
		delete(f.dependencies, name)
	})
}
