package views

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/slices"
)

type editView struct {
	updateInputView   func()
	editViewContainer *fyne.Container
	changeVariable    func(*formulaInfo)
}

func newEditView(variables map[string]*formulaInfo, parentWindow fyne.Window) *editView {
	// Create the editor
	variableEditor := widget.NewMultiLineEntry()
	variableEditor.SetPlaceHolder("Formula")
	editorVariable := ""

	// Formula inputs selection
	selectedInput := ""
	inputVariableSelect := widget.NewSelect([]string{}, func(s string) {
		selectedInput = s
		log.Println(selectedInput)
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
		log.Println("Updating input display for:", editorVariable)
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
	inputView := container.NewBorder(nil, inputDisplay, nil, addInputButton, inputVariableSelect)

	// Add a button to change to edit mode
	nameLabel := widget.NewLabel(editorVariable)
	editNameButton := widget.NewButton("Rename", func() {})
	editNameButton.OnTapped = func() {

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

			// Check for valid characters
			letters := `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
			letters += `abcdefghijklmnopqrstuvwxyz`
			validCharacters := letters
			validCharacters += `0123456789`
			validCharacters += `_`
			for index, character := range input {
				characterString := string(character)
				if index == 1 && !strings.Contains(letters, characterString) {
					log.Println("Invalid variable name")
					return errors.New(`"` + characterString + "\" is not a valid 1st character")
				}
				if !strings.Contains(validCharacters, characterString) {
					log.Println("Invalid variable name")
					return errors.New(`"` + characterString + "\" is not a valid character")
				}
			}
			return nil
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
				fmt.Printf("%s dependencies: %v\n", dependentName, variables[dependentName].dependencies)
			}
		}, parentWindow)
	}

	// Build the view
	nameView := container.NewBorder(nil, nil, nil, editNameButton, container.New(layout.NewCenterLayout(), nameLabel))
	return &editView{
		editViewContainer: container.NewBorder(
			container.NewBorder(nameView, nil, nil, nil, inputView),
			nil, nil, nil, variableEditor),
		updateInputView: func() {
			updateInputSelect()
			updateInputDisplay()
		},
		changeVariable: func(variable *formulaInfo) {
			// Get the variable name
			name, err := variable.name.Get()
			checkErrFatal("Failed to get variable name:", err)

			// Change out the editor components
			editorVariable = name
			nameLabel.Bind(variable.name)
			variableEditor.Bind(variable.code)
			updateInputSelect()
			updateInputDisplay()
		},
	}
}
