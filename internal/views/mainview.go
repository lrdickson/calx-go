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

	// Display the output
	outputListBind := binding.NewStringList()
	outputListDisplay := widget.NewListWithData(outputListBind,
		func() fyne.CanvasObject {
			return newVariableView2()
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*VariableView).Bind(i.(binding.String))
		})
	newVariableButton := widget.NewButton("New", func() { outputListBind.Append("Hello") })

	// Display data from the variable
	variableDisplay1 := widget.NewLabel("")

	// Create the editor
	editorCodeBind := binding.BindString(&mainViewModel.EditorCode)
	variableEditor := widget.NewMultiLineEntry()
	variableEditor.Bind(editorCodeBind)
	variableEditor.SetPlaceHolder("Enter text...")

	// Run variable code button
	runButton := widget.NewButton("Run", func() {
		variableDisplay1.SetText(mainViewModel.RunCode())
	})

	// Put everything together
	content := container.NewBorder(
		nil, nil,
		container.NewBorder(variableDisplay1, newVariableButton, nil, nil, outputListDisplay),
		nil,
		container.NewBorder(nil, runButton, nil, nil, variableEditor))

	return content
}
