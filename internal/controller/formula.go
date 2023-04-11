package controller

type Formula struct {
	BaseVariable
	code         string
	dependencies map[string]bool
	dependents   map[string]bool
}

func (f Formula) Code() string {
	return f.code
}

func AddFormula(c *Controller) {
	var formula Variable = &Formula{}
	c.AddVariable(c.UniqueName(), &formula)
}

func (f *Formula) AddDependency(c *Controller, name string) {
	f.dependencies[name] = true
	c.AddListener(RenameVarEvent, name, func(dependencyName string) {
		f.dependencies[dependencyName] = true
		delete(f.dependencies, name)
	})
}
