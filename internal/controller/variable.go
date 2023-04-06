package controller

type Variable interface {
	Name() string
	Data() any
}

type baseVariable struct {
	name string
	data any
}

func (v baseVariable) Name() string {
	return v.name
}

func (v baseVariable) Data() any {
	return v.data
}

type Formula struct {
	baseVariable
	code string
}

func (f Formula) Code() string {
	return f.code
}
