package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/lrdickson/ssgo/internal/viewmodels"
)

func NewMainView() *fyne.Container {
	mainViewModel := viewmodels.NewMainViewModel()

	// Display the variable
	variableListBind := binding.BindStringList(&mainViewModel.VariableList)
	variableListDisplay := widget.NewListWithData(variableListBind,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	newVariableButton := widget.NewButton("New", func() { variableListBind.Append("Hello") })

	// Display data from the variable
	variableDisplay := widget.NewLabel("")

	// Create the editor
	editorCodeBind := binding.BindString(&mainViewModel.EditorCode)
	variableEditor := widget.NewMultiLineEntry()
	variableEditor.Bind(editorCodeBind)
	variableEditor.SetPlaceHolder("Enter text...")

	// Run variable code button
	runButton := widget.NewButton("Run", func() {
		variableDisplay.SetText(mainViewModel.RunCode())
	})

	// Put everything together
	content := container.NewBorder(
		nil, nil,
		container.NewBorder(
			nil, nil,
			container.NewBorder(nil, newVariableButton, nil, nil, variableListDisplay),
			nil, variableDisplay),
		nil,
		container.NewBorder(nil, runButton, nil, nil, variableEditor))

	return content
}
