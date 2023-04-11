package controller

type BaseVariable struct {
	data any
}

func (v *BaseVariable) Data() any {
	return v.data
}

func (v *BaseVariable) Close() {}
