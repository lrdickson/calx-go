package views

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/lrdickson/ssgo/internal/viewmodels"
)

type formulaInfo struct {
	code   binding.String
	name   binding.String
	output binding.String
}

func checkErrFatal(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func getVariable(variables binding.UntypedList, id widget.ListItemID) formulaInfo {
	variablesInterface, err := variables.Get()
	checkErrFatal("Failed to get variable interface array:", err)
	return variablesInterface[id].(formulaInfo)
}

func NewMainView() *fyne.Container {

	// Create the editor
	variableEditor := widget.NewMultiLineEntry()
	variableEditor.SetPlaceHolder("Enter text...")

	// Display the output
	variables := binding.NewUntypedList()
	variableList := widget.NewListWithData(
		variables,
		func() fyne.CanvasObject {
			nameDisplay := widget.NewLabel("AAAAAAAAAAAAAAAAAAAAAAA")
			nameEditor := widget.NewEntry()
			nameEditor.Hide()
			editNameButton := widget.NewButton("Edit", func() {})
			editNameButton.OnTapped = func() {
				if nameDisplay.Visible() {
					nameDisplay.Hide()
					nameEditor.Show()
					editNameButton.SetText("Update")
				} else {
					nameDisplay.Show()
					nameEditor.Hide()
					editNameButton.SetText("Edit")
				}
			}
			name := container.NewBorder(nil, nil, nil, editNameButton, container.NewMax(nameDisplay, nameEditor))
			//name := container.NewBorder(
			//// The width of the variable pane can be controlled by the length of this label
			//widget.NewLabel("AAAAAAAAAAAAAAAAAAAAAAA"),
			//nil, nil, nil,
			//widget.NewEntry())
			output := widget.NewLabel("Output")
			return container.NewBorder(name, nil, nil, nil, output)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			// Get the variable
			v, err := item.(binding.Untyped).Get()
			checkErrFatal("Failed to get variable data:", err)
			variable := v.(formulaInfo)

			// Set the output
			output := obj.(*fyne.Container).Objects[0].(*widget.Label)
			output.Bind(variable.output)
			output.Refresh()

			// Set the name
			name := obj.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label)
			name.Bind(variable.name)
			nameEntry := obj.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Entry)
			nameEntry.Bind(variable.name)
			name.Refresh()
		})

	// Create a new variable
	variableCount := 1
	newVariableButton := widget.NewButton("New", func() {
		// Add the variable name
		name := binding.NewString()
		name.Set("NewVariable" + strconv.Itoa(variableCount))
		variableCount++

		// Build the variable
		code := binding.NewString()
		code.Set("")
		output := binding.NewString()
		output.Set("")
		newVariable := formulaInfo{code, name, output}
		variables.Append(newVariable)
		variableList.Refresh()
	})

	// Edit the code of the selected variable
	variableList.OnSelected = func(id widget.ListItemID) {
		// Assign the code to the editor
		code := getVariable(variables, id).code
		variableEditor.Bind(code)
	}

	// Run variable code button
	mainViewModel := viewmodels.NewMainViewModel()
	runButton := widget.NewButton("Run", func() {
		ivariables, err := variables.Get()
		checkErrFatal("Failed to get variable interface array:", err)
		for _, ivariable := range ivariables {
			variable := ivariable.(formulaInfo)
			code, err := variable.code.Get()
			checkErrFatal("Failed to get formula code:", err)
			mainViewModel.EditorCode = code
			variable.output.Set(mainViewModel.RunCode())
		}
	})

	// Put everything together
	content := container.NewBorder(
		nil, nil,
		container.NewBorder(nil, newVariableButton, nil, nil, variableList),
		nil,
		container.NewBorder(nil, runButton, nil, nil, variableEditor))

	return content
}
