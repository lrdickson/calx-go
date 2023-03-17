package views

import (
	"fyne.io/fyne/v2/data/binding"
)

type variableListView struct {
	namesBinding binding.StringList
}

//func newVariableListView() *variableListView {

//view := variableListView{binding.NewStringList()}

//list := widget.NewList()

//// Display the output
//outputListBind := binding.NewStringList()
//outputListDisplay := widget.NewListWithData(outputListBind,
//func() fyne.CanvasObject {
//return newVariableView2()
//},
//func(i binding.DataItem, o fyne.CanvasObject) {
//o.(*VariableView).Bind(i.(binding.String))
//})
//newVariableButton := widget.NewButton("New", func() { outputListBind.Append("Hello") })

//}
