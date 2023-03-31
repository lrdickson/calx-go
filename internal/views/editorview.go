package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/slices"
)

type editView struct {
	updateInputView   func()
	editViewContainer *fyne.Container
	changeVariable    func(string, binding.String)
}

func newEditView(variables map[string]*formulaInfo) *editView {
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

	// Build the view
	inputView := container.NewBorder(nil, inputDisplay, nil, addInputButton, inputVariableSelect)
	return &editView{
		editViewContainer: container.NewBorder(inputView, nil, nil, nil, variableEditor),
		updateInputView: func() {
			updateInputSelect()
			updateInputDisplay()
		},
		changeVariable: func(name string, code binding.String) {
			editorVariable = name
			variableEditor.Bind(code)
			updateInputSelect()
			updateInputDisplay()
		},
	}
}
