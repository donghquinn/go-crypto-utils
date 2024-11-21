package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"org.donghyuns.com/secure/keygen/gui"
)

type testInterface struct {
	Email string `json:"email"`
}

func main() {
	// Initialize the Fyne app
	myApp := app.New()
	myWindow := myApp.NewWindow("Encryption Generator")
	myWindow.Resize(fyne.NewSize(800, 800))

	// Build the GUI and set it as the window content
	content := gui.BuildGUI(myApp)
	myWindow.SetContent(content)

	// Show and run the application
	myWindow.ShowAndRun()
}
