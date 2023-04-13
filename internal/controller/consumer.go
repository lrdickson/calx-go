package controller

type Consumer interface {
	Eval(map[string]any)
}
