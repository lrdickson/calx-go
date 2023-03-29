package views

import (
	"errors"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func newVariableDisplayView(variables map[string]*formulaInfo, updateInputView func(), parent fyne.Window) (binding.UntypedList, *widget.List) {

	// Display the output
	displayVariables := binding.NewUntypedList()
	displayVariablesView := widget.NewListWithData(
		displayVariables,
		func() fyne.CanvasObject {
			// Add name the elements
			nameDisplay := widget.NewLabel("")
			editNameButton := widget.NewButton("Rename", func() {})

			// Add a button to change to edit mode
			editNameButton.OnTapped = func() {

				// Create the name editor form item
				nameEditor := widget.NewEntry()
				nameEditor.SetText(nameDisplay.Text)
				oldName := nameDisplay.Text
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
							fmt.Println("Invalid variable name")
							return errors.New(`"` + characterString + "\" is not a valid 1st character")
						}
						if !strings.Contains(validCharacters, characterString) {
							fmt.Println("Invalid variable name")
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
					updateInputView()
				}, parent)
			}
			name := container.NewBorder(nil, nil, nil, editNameButton, nameDisplay)
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

			// Set the name
			nameLabel := obj.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Label)
			nameLabel.Bind(variable.name)
		})

	return displayVariables, displayVariablesView
}
