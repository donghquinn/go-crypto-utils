package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"org.donghyuns.com/secure/keygen/gui"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Secure Helper")
	myWindow.Resize(fyne.NewSize(800, 600))

	// Create Tabs
	tabs := container.NewAppTabs(
		gui.NewEncryptTab(myApp, myWindow),
		gui.NewDecryptTab(myApp, myWindow),
		gui.NewEncodeTab(myApp, myWindow),
		gui.NewDecodeTab(myApp, myWindow),
		// gui.NewDecodeBase64Tab(myApp, myWindow),
		gui.NewKeyGenTab(myApp, myWindow),
		gui.NewUUIDTab(myApp, myWindow),
		gui.NewRandomStringTab(myApp, myWindow),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
