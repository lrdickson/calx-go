package kernel

import "github.com/lrdickson/calx/internal/controller"

type Kernel interface {
}

type FormulaEngine interface {
	controller.Consumer
	controller.Producer
	SetCode(string)
	Code() string
}
