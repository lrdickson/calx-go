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
	"golang.org/x/exp/slices"
)

type formulaInfo struct {
	code         binding.String
	name         binding.String
	output       binding.String
	dependencies map[string]*formulaInfo
	dependents   map[string]*formulaInfo
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
		inputVariableSelect.Options = variableSelectList
		if !slices.Contains(variableSelectList, inputVariableSelect.Selected) {
			inputVariableSelect.ClearSelected()
		}
	}
	inputDisplay := container.NewHScroll(container.NewHBox())
	inputDisplay.Hide()

	// Edit the code of the selected variable
	var updateInputDisplay func()
	updateInputDisplay = func() {
		fmt.Println("Updating input display for:", editorVariable)
		if _, exists := variables[editorVariable]; !exists {
			return
		}
		if len(variables[editorVariable].dependencies) == 0 {
			inputDisplay.Hide()
			return
		}

		// Create a list of buttons to display
		inputArray := make([]fyne.CanvasObject, 0, len(variables[editorVariable].dependencies))
		for inputVariable := range variables[editorVariable].dependencies {
			fmt.Printf("Adding %s to input display\n", inputVariable)

			// Make a copy so that the variable being deleted does change as the value of inputVariable changes
			buttonVariable := inputVariable
			inputArray = append(inputArray, widget.NewButton(inputVariable+" X", func() {
				delete(variables[editorVariable].dependencies, buttonVariable)
				updateInputDisplay()
			}))
		}
		inputDisplay.Content = container.NewHBox(inputArray...)
		inputDisplay.Refresh()
		inputDisplay.Show()
	}
	updateInputView := func() {
		updateInputSelect()
		updateInputDisplay()
	}

	// Display the output
	displayVariables, displayVariablesView := newVariableDisplayView(variables, updateInputView, parent)
	displayVariablesView.OnSelected = func(id widget.ListItemID) {
		// Get the variable
		variable := getVariable(displayVariables, id)
		name, err := variable.name.Get()
		checkErrFatal("Failed to get variable name:", err)
		editorVariable = name

		// Update the editor
		variableEditor.Bind(variable.code)
		updateInputView()
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
		variables[editorVariable].dependencies[selectedInput] = variables[selectedInput]
		variables[selectedInput].dependents[editorVariable] = variables[editorVariable]
		updateInputDisplay()
	})
	inputView := container.NewBorder(nil, inputDisplay, nil, addInputButton, inputVariableSelect)

	// Create a new variable
	variableCount := 1
	newVariableButton := widget.NewButton("New", func() {
		// Add the variable nameDisplay
		name := ""
		for {
			name = "var" + strconv.Itoa(variableCount)
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
		newVariable := formulaInfo{code, nameDisplay, output, make(map[string]*formulaInfo), make(map[string]*formulaInfo)}
		displayVariables.Append(newVariable)
		variables[name] = &newVariable
		updateInputSelect()
	})

	// Run variable code button
	goKernel := kernel.NewKernel()
	runButton := widget.NewButton("Run", func() {
		input := make(map[string]*kernel.Formula)
		for name := range variables {
			code, err := variables[name].code.Get()
			checkErrFatal("Failed to get formula code:", err)
			dependencies := make([]string, 0, len(variables[name].dependencies))
			for dependencyName := range variables[name].dependencies {
				dependencies = append(dependencies, dependencyName)
			}
			input[name] = &kernel.Formula{Code: code, Dependencies: dependencies}
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
