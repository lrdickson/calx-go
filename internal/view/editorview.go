package view

import (
	"errors"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lrdickson/calx/internal/controller"
	"github.com/lrdickson/calx/internal/variable"
	"golang.org/x/exp/slices"
)

type editView struct {
	editViewContainer *fyne.Container
	updateEditorView  func(*formulaInfo)
}

func updateRenameFunction(editorVariable string, variables map[string]*formulaInfo, parentWindow fyne.Window) func() {
	return func() {

		// Create the name editor form item
		nameEditor := widget.NewEntry()
		nameEditor.SetText(editorVariable)
		oldName := editorVariable
		nameEditor.Validator = func(input string) error {
			// Check if the name is taken
			_, taken := variables[input]
			if oldName != input && taken {
				return errors.New(input + " is already taken")
			}

			// Check if the name is valid
			return variable.NameValid(input)
		}
		nameItem := &widget.FormItem{
			Widget: nameEditor,
		}

		// Show the form
		items := []*widget.FormItem{nameItem}
		dialog.ShowForm("Update Formula Name", "Submit", "Cancel", items, func(confirm bool) {
			// Do nothing if cancelled
			if !confirm {
				return
			}

			// Check if the name changed
			newName := nameEditor.Text
			if newName == oldName {
				return
			}

			// Update the variable
			variables[oldName].name.Set(newName)
			variables[newName] = variables[oldName]
			delete(variables, oldName)
			for dependentName := range variables[newName].dependents {
				variables[dependentName].dependencies[newName] = variables[newName]
				delete(variables[dependentName].dependencies, oldName)
				log.Printf("%s dependencies: %v\n", dependentName, variables[dependentName].dependencies)
			}
		}, parentWindow)
	}
}

func newInputView(editorVariable string, control *controller.Controller) (*fyne.Container, func(string)) {
	// Formula inputs selection
	selectedInput := ""
	inputVariableSelect := widget.NewSelect([]string{}, func(s string) {
		selectedInput = s
		log.Println("Selected input:", selectedInput)
	})
	updateInputSelect := func(editorVariable string) {
		variableSelectList := make([]string, 0, control.VariableCount())
		control.IterVariables(func(name string, v *variable.Variable) bool {
			if name != editorVariable {
				variableSelectList = append(variableSelectList, name)
			}
			return true
		})
		inputVariableSelect.Options = variableSelectList
		if !slices.Contains(variableSelectList, inputVariableSelect.Selected) {
			inputVariableSelect.ClearSelected()
		}
	}
	inputDisplay := container.NewHScroll(container.NewHBox())
	inputDisplay.Hide()

	// Edit the code of the selected variable
	var updateInputDisplay func(editorVariable string)
	updateInputDisplay = func(editorVariable string) {
		log.Println("Updating input display for:", editorVariable)
		if v := control.Variables(editorVariable); v == nil {
			return
		}
		if len(control.Variables(editorVariable).dependencies) == 0 {
			inputDisplay.Hide()
			return
		}

		// Create a list of buttons to display
		inputArray := make([]fyne.CanvasObject, 0, len(control.Variables(editorVariable).dependencies))
		for inputVariable := range control.Variables(editorVariable).dependencies {
			log.Printf("Adding %s to input display\n", inputVariable)

			// Make a copy so that the variable being deleted does change as the value of inputVariable changes
			buttonVariable := inputVariable
			inputArray = append(inputArray, widget.NewButton(inputVariable+" X", func() {
				delete(control.Variables(editorVariable).dependencies, buttonVariable)
				updateInputDisplay(editorVariable)
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
		updateInputDisplay(editorVariable)
	})
	inputView := container.NewBorder(nil, inputDisplay, nil, addInputButton, inputVariableSelect)
	updateInputView := func(editorVariable string) {
		updateInputSelect(editorVariable)
		updateInputDisplay(editorVariable)
		addInputButton.OnTapped = func() {
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
			updateInputDisplay(editorVariable)
		}
	}
	return inputView, updateInputView
}

func updateDeleteFunction() func() {
	return func() {}
}

func newEditView(control *controller.Controller, parentWindow fyne.Window) *editView {
	// Create the editor
	variableEditor := widget.NewMultiLineEntry()
	variableEditor.SetPlaceHolder("Formula")
	editorVariable := ""

	// Create the input view
	inputView, updateInputView := newInputView(editorVariable, control)

	// Add the delete button
	deleteButton := widget.NewButton("Delete", nil)

	// Add the name label
	editNameButton := widget.NewButton("Rename", nil)
	nameLabel := widget.NewLabel(editorVariable)
	nameView := container.NewBorder(nil, nil, nil, container.NewHBox(editNameButton, deleteButton),
		container.New(layout.NewCenterLayout(), nameLabel))

	// Build the view
	var previousVariable *formulaInfo
	return &editView{
		editViewContainer: container.NewBorder(
			container.NewBorder(nameView, nil, nil, nil, inputView),
			nil, nil, nil, variableEditor),
		updateEditorView: func(variable *formulaInfo) {
			// Return if variable doesn't exist
			if variable == nil {
				return
			}

			// Get the variable name
			name, err := variable.name.Get()
			checkErrFatal("Failed to get variable name:", err)

			if editorVariable != name {
				editorVariable = name
				editNameButton.OnTapped = updateRenameFunction(name, variables, parentWindow)
				deleteButton.OnTapped = updateDeleteFunction()
			}

			if previousVariable != variable {
				previousVariable = variable
				nameLabel.Bind(variable.name)
				variableEditor.Bind(variable.code)
			}

			// This will probably change every time
			updateInputView(name)
		},
	}
}
