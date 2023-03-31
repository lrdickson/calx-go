package views

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/lrdickson/ssgo/internal/kernel"
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

func RunGui() {
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("SSGO")

	variables := make(map[string]*formulaInfo)

	mainEditView := newEditView(variables)

	// Display the output
	displayVariables, displayVariablesView := newVariableDisplayView(variables, mainEditView, mainWindow)

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
		mainEditView.updateInputView()
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
		container.NewBorder(nil, runButton, nil, nil, mainEditView.editViewContainer))

	mainWindow.SetContent(content)
	mainWindow.Resize(fyne.NewSize(480, 360))
	mainWindow.ShowAndRun()
}
