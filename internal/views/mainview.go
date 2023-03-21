package views

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/lrdickson/ssgo/internal/kernel"
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

func NewMainView() *container.Split {

	// Create the editor
	variableEditor := widget.NewMultiLineEntry()
	variableEditor.SetPlaceHolder("Enter text...")

	// Display the output
	variables := binding.NewUntypedList()
	variableList := widget.NewListWithData(
		variables,
		func() fyne.CanvasObject {
			nameDisplay := widget.NewLabel("")
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
			name := obj.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container)
			nameLabel := name.Objects[0].(*widget.Label)
			nameLabel.Bind(variable.name)
			nameEntry := name.Objects[1].(*widget.Entry)
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
	goKernel := kernel.NewKernel()
	runButton := widget.NewButton("Run", func() {
		ivariables, err := variables.Get()
		checkErrFatal("Failed to get variable interface array:", err)
		input := make(map[string]string)
		for _, ivariable := range ivariables {
			variable := ivariable.(formulaInfo)
			code, err := variable.code.Get()
			checkErrFatal("Failed to get formula code:", err)
			name, err := variable.name.Get()
			checkErrFatal("Failed to get formula name:", err)
			input[name] = code
		}
		output := goKernel.Update(input)
		for _, ivariable := range ivariables {
			variable := ivariable.(formulaInfo)
			name, err := variable.name.Get()
			checkErrFatal("Failed to get formula name:", err)
			variable.output.Set(output[name])
		}
	})

	// Put everything together
	content := container.NewHSplit(
		container.NewBorder(nil, newVariableButton, nil, nil, variableList),
		container.NewBorder(nil, runButton, nil, nil, variableEditor))

	return content
}
