package controller

type Variable interface {
	Data() any
}

type baseVariable struct {
	data any
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
