package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func newVariableDisplayView(variables map[string]*formulaInfo) (binding.UntypedList, *widget.List) {

	// Display the output
	displayVariables := binding.NewUntypedList()
	displayVariablesView := widget.NewListWithData(
		displayVariables,
		func() fyne.CanvasObject {
			// Add name the elements
			nameDisplay := widget.NewLabel("")
			output := widget.NewLabel("Output")
			return container.NewBorder(nameDisplay, nil, nil, nil, output)
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
			nameLabel := obj.(*fyne.Container).Objects[1].(*widget.Label)
			nameLabel.Bind(variable.name)
		})

	return displayVariables, displayVariablesView
}
