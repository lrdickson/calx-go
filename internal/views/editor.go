package views

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewEditor() *widget.List {
	return widget.NewList(
		func() int {
			return 1
		},
		func() fyne.CanvasObject {
			return container.NewHBox(canvas.NewText("", color.White))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {})
}
