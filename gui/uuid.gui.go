package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"org.donghyuns.com/secure/keygen/biz"
)

// NewUUIDTab creates the UUID tab
func NewUUIDTab(app fyne.App, window fyne.Window) *container.TabItem {
	uuidEntry := widget.NewEntry()
	uuidEntry.SetPlaceHolder("Generated UUID will appear here...")

	generateBtn := widget.NewButton("Generate UUID", func() {
		newUUID := biz.CreateUuid()
		uuidEntry.SetText(newUUID)
	})

	copyBtn := widget.NewButton("Copy", func() {
		if uuidEntry.Text == "" {
			dialog.ShowInformation("No UUID", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(uuidEntry.Text)
		dialog.ShowInformation("Copied", "UUID copied to clipboard.", window)
	})

	content := container.NewVBox(
		widget.NewLabel("UUID:"),
		uuidEntry,
		container.NewHBox(generateBtn, copyBtn),
	)

	return container.NewTabItem("UUID", content)
}
