package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/lrdickson/ssgo/internal/views"
)

func main() {
	// Start the GUI
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("SSGO")
	mainWindow.SetContent(views.NewMainView())
	mainWindow.Resize(fyne.NewSize(480, 360))
	mainWindow.ShowAndRun()
}
