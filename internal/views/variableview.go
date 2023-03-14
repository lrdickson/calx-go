package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func newVariableView() fyne.CanvasObject {
	return widget.NewLabel("template")
}

func newVariableView2() *VariableView {
	variableView := &VariableView{}
	variableView.name = widget.NewLabel("")
	variableView.view = container.NewBorder(nil, nil, nil, nil, variableView.name)
	return variableView
}

type VariableView struct {
	view *fyne.Container
	name *widget.Label
}

func (v *VariableView) MinSize() fyne.Size {
	return v.view.MinSize()
}

func (v *VariableView) Move(p fyne.Position) {
	v.view.Move(p)
}

func (v *VariableView) Position() fyne.Position {
	return v.view.Position()
}

func (v *VariableView) Resize(s fyne.Size) {
	v.view.Resize(s)
}

func (v *VariableView) Size() fyne.Size {
	return v.view.Size()
}

func (v *VariableView) Hide() {
	v.view.Hide()
}

func (v *VariableView) Visible() bool {
	return v.view.Visible()
}

func (v *VariableView) Show() {
	v.view.Show()
}

func (v *VariableView) Refresh() {
	v.view.Refresh()
	v.name.Refresh()
}

func (v *VariableView) Bind(data binding.String) {
	v.name.Bind(data)
}
