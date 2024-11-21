package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"org.donghyuns.com/secure/keygen/gui"
)

type testInterface struct {
	Email string `json:"email"`
}

func main() {
	// Initialize the Fyne app
	myApp := app.New()
	myWindow := myApp.NewWindow("Encryption Generator")
	myWindow.Resize(fyne.NewSize(500, 500))
	// myWindow.SetIcon(theme.DocumentIcon()) // Optional: Set an icon for the window

	// Build the GUI and set it as the window content
	content := gui.BuildGUI(myApp)
	myWindow.SetContent(content)

	// Show and run the application
	myWindow.ShowAndRun()
}
