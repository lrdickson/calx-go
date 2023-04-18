package gokernel

import (
	"github.com/lrdickson/calx/internal/controller"
)

type FormulaEngine struct {
	controller.BaseObjectEngine
	code            string
	output          any
	onOutputChanged func()
}

func NewFormula(c *controller.Controller) *controller.Object {
	var formula controller.ObjectEngine = &FormulaEngine{}
	return c.NewObject(c.UniqueName(), &formula)
}

func (f *FormulaEngine) Code() string {
	return f.code
}

func (f *FormulaEngine) SetCode(code string) {
	f.code = code
}

func (f *FormulaEngine) Consume(any) error {
	return nil
}

func (f *FormulaEngine) Output() (any, error) {
	return f.output, nil
}

func (f *FormulaEngine) SetOnOutputChanged(changed func()) {
	f.onOutputChanged = changed
}
