package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/lrdickson/ssgo/internal/views"
)

func main() {
	// Start the GUI
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("SSGO")
	mainWindow.SetContent(views.NewMainView())
	mainWindow.ShowAndRun()
}
