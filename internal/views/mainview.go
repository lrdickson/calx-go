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
	inputs map[string]*formulaInfo
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
	editorVariable := ""

	// Formula inputs selection
	selectedInput := ""
	inputVariableSelect := widget.NewSelect([]string{}, func(s string) {
		selectedInput = s
		fmt.Println(selectedInput)
	})
	updateInputSelect := func() {
		variableSelectList := make([]string, 0, len(variables))
		for name := range variables {
			if name != editorVariable {
				variableSelectList = append(variableSelectList, name)
			}
		}
		inputVariableSelect.SetSelected("")
		inputVariableSelect.Options = variableSelectList
	}
	inputDisplay := container.NewHScroll(container.NewHBox())
	inputDisplay.Hide()

	// Display the output
	displayVariables, displayVariablesView := newVariableDisplayView(variables, updateInputSelect, parent)

	// Edit the code of the selected variable
	var updateInputDisplay func()
	updateInputDisplay = func() {
		if _, exists := variables[editorVariable]; !exists {
			return
		}
		if len(variables[editorVariable].inputs) == 0 {
			inputDisplay.Hide()
			return
		}

		// Create a list of buttons to display
		inputArray := make([]fyne.CanvasObject, 0, len(variables[editorVariable].inputs))
		for inputVariable := range variables[editorVariable].inputs {
			inputArray = append(inputArray, widget.NewButton(inputVariable+" X", func() {
				delete(variables[editorVariable].inputs, inputVariable)
				updateInputDisplay()
			}))
		}
		inputDisplay.Content = container.NewHBox(inputArray...)
		inputDisplay.Show()
	}
	displayVariablesView.OnSelected = func(id widget.ListItemID) {
		// Get the variable
		variable := getVariable(displayVariables, id)
		name, err := variable.name.Get()
		checkErrFatal("Failed to get variable name:", err)
		editorVariable = name

		// Update the editor
		variableEditor.Bind(variable.code)
		updateInputSelect()
		updateInputDisplay()
	}

	// Button to add selected inputs to a formula
	addInputButton := widget.NewButton("Add Input", func() {
		if selectedInput == editorVariable {
			return
		}
		if _, exists := variables[selectedInput]; !exists {
			return
		}
		if _, exists := variables[editorVariable]; !exists {
			return
		}
		variables[editorVariable].inputs[selectedInput] = variables[selectedInput]
		updateInputDisplay()
	})
	inputView := container.NewVBox(
		container.NewHBox(inputVariableSelect, addInputButton),
		inputDisplay)

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
		newVariable := formulaInfo{code, nameDisplay, output, make(map[string]*formulaInfo)}
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
