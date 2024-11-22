package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"org.donghyuns.com/secure/keygen/gui"
)

type testInterface struct {
	Email string `json:"email"`
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("AES Key Generator")
	myWindow.Resize(fyne.NewSize(800, 600))

	// Create Tabs
	tabs := container.NewAppTabs(
		gui.NewEncryptTab(myApp, myWindow),
		gui.NewDecryptTab(myApp, myWindow),
		gui.NewBase64Tab(myApp, myWindow),
		// gui.NewDecodeBase64Tab(myApp, myWindow),
		gui.NewKeyGenTab(myApp, myWindow),
		gui.NewUUIDTab(myApp, myWindow),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
