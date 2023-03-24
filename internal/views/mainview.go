package views

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
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

	// Create the editor
	variableEditor := widget.NewMultiLineEntry()
	variableEditor.SetPlaceHolder("Formula")

	// Display the output
	variables := make(map[string]*formulaInfo)
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
	})

	// Edit the code of the selected variable
	displayVariablesView.OnSelected = func(id widget.ListItemID) {
		// Assign the code to the editor
		code := getVariable(displayVariables, id).code
		variableEditor.Bind(code)
	}

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
		container.NewBorder(nil, runButton, nil, nil, variableEditor))

	return content
}
