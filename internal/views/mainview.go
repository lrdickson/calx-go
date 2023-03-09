package views

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/lrdickson/ssgo/internal/runner"
	"github.com/lrdickson/ssgo/internal/viewmodels"
)

func NewMainView() *fyne.Container {
	mainViewModel := viewmodels.NewMainViewModel()

	// Display the variable
	//variableList = binding.NewStringList()
	variableDisplay := widget.NewListWithData(mainViewModel.VariableList,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	newVariableButton := widget.NewButton("New", func() { mainViewModel.VariableList.Append("Hello") })

	variableDisplayOld := canvas.NewText("", color.White)

	// Create the editor
	variableEditor := widget.NewMultiLineEntry()
	variableEditor.SetPlaceHolder("Enter text...")

	// Run variable code button
	runButton := widget.NewButton("Run", func() {
		log.Println("Running:", variableEditor.Text)
		display, err := runner.RunGo(variableEditor.Text)
		if err != nil {
			display = "Err"
		}
		variableDisplayOld.Text = display
		variableDisplayOld.Color = color.Black
		variableDisplayOld.Refresh()
	})

	// Put everything together
	content := container.NewBorder(nil, nil,
		container.NewBorder(variableDisplayOld, newVariableButton, nil, nil, variableDisplay),
		nil,
		container.NewBorder(nil, runButton, nil, nil, variableEditor))

	return content
}
