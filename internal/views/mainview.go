package views

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/lrdickson/ssgo/internal/viewmodels"
)

type variableInfo struct {
	code   string
	name   string
	output string
}

func checkErrFatal(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func NewMainView() *fyne.Container {

	// Display the output
	variables := binding.NewUntypedList()
	variableList := widget.NewListWithData(
		variables,
		func() fyne.CanvasObject {
			name := widget.NewLabel("Template")
			output := widget.NewLabel("Output")
			return container.NewBorder(name, nil, nil, nil, output)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			// Get the variable
			v, err := item.(binding.Untyped).Get()
			checkErrFatal("Failed to get variable data:", err)
			variable := v.(variableInfo)

			// Set the output
			output := obj.(*fyne.Container).Objects[0].(*widget.Label)
			output.SetText(variable.output)
			output.Refresh()

			// Set the name
			name := obj.(*fyne.Container).Objects[1].(*widget.Label)
			name.SetText(variable.name)
			name.Refresh()
		})
	newVariableButton := widget.NewButton("New", func() {
		// Create a new variable
		newVariable := variableInfo{"", "NewVariable", ""}
		variables.Append(newVariable)
		variableList.Refresh()
	})

	// Display data from the variable
	variableDisplay1 := widget.NewLabel("")

	// Create the editor
	mainViewModel := viewmodels.NewMainViewModel()
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
		container.NewBorder(variableDisplay1, newVariableButton, nil, nil, variableList),
		nil,
		container.NewBorder(nil, runButton, nil, nil, variableEditor))

	return content
}
