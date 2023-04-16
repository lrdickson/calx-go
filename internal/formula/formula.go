package formula

import "github.com/lrdickson/calx/internal/controller"

type Formula struct {
	controller.BaseObjectEngine
	code            string
	output          any
	onOutputChanged func()
}

func (f *Formula) Code() string {
	return f.code
}

func NewFormula(c *controller.Controller) *controller.Object {
	var formula controller.ObjectEngine = &Formula{}
	return c.NewObject(c.UniqueName(), &formula)
}

func Kernel() string {
	return "Go"
}

func Consume(any) error {
	return nil
}

func (f *Formula) Output() (any, error) {
	return f.output, nil
}

func (f *Formula) SetOnOutputChanged(changed func()) {
	f.onOutputChanged = changed
}
