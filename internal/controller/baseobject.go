package controller

type BaseObject struct {
	name string
}

func (b *BaseObject) Close() {}

func (b *BaseObject) Name() string { return b.name }
