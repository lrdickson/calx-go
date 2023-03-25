package views

import (
	"fmt"
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

func NewMainView(parent fyne.Window) *container.Split {

	variables := make(map[string]*formulaInfo)

	// Create the editor
	variableEditor := widget.NewMultiLineEntry()
	variableEditor.SetPlaceHolder("Formula")

	// Formula inputs
	inputSelect := ""
	inputSelected := func(s string) {
		inputSelect = s
		fmt.Println(inputSelect)
	}
	variableSelect := widget.NewSelect([]string{}, inputSelected)
	updateInputSelect := func() {
		variableSelectList := make([]string, 0, len(variables))
		for name := range variables {
			variableSelectList = append(variableSelectList, name)
		}
		variableSelect.Options = variableSelectList
	}

	// Button to add selected inputs to a formula
	inputView := container.NewVBox(variableSelect)

	// Display the output
	displayVariables, displayVariablesView := newVariableDisplayView(variables, updateInputSelect, parent)

	// Edit the code of the selected variable
	displayVariablesView.OnSelected = func(id widget.ListItemID) {
		// Assign the code to the editor
		code := getVariable(displayVariables, id).code
		variableEditor.Bind(code)
	}

	// Create a new variable
	variableCount := 1
	newVariableButton := widget.NewButton("New", func() {
		// Add the variable nameDisplay
		name := ""
		for {
			name = "NewVariable" + strconv.Itoa(variableCount)
			variableCount++
			if _, taken := variables[name]; !taken {
				break
			}
		}
		nameDisplay := binding.NewString()
		nameDisplay.Set(name)

		// Build the variable
		code := binding.NewString()
		output := binding.NewString()
		newVariable := formulaInfo{code, nameDisplay, output}
		displayVariables.Append(newVariable)
		variables[name] = &newVariable
		updateInputSelect()
	})

	// Run variable code button
	goKernel := kernel.NewKernel()
	runButton := widget.NewButton("Run", func() {
		input := make(map[string]string)
		for name := range variables {
			code, err := variables[name].code.Get()
			checkErrFatal("Failed to get formula code:", err)
			input[name] = code
		}
		output := goKernel.Update(input)
		for name, variable := range variables {
			variable.output.Set(output[name])
		}
	})

	// Put everything together
	content := container.NewHSplit(
		container.NewBorder(nil, newVariableButton, nil, nil, displayVariablesView),
		container.NewBorder(inputView, runButton, nil, nil, variableEditor))

	return content
}
