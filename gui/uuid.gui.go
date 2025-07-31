package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"org.donghyuns.com/secure/keygen/biz"
)

// NewUUIDTab creates the UUID tab
func NewUUIDTab(app fyne.App, window fyne.Window) *container.TabItem {
	ui := NewUIComponents(app, window)
	uuidEntry := widget.NewEntry()
	uuidEntry.SetPlaceHolder("Generated UUID will appear here...")

	generateBtn := ui.CreateProcessButton("Generate UUID", func() {
		newUUID := biz.CreateUuid()
		uuidEntry.SetText(newUUID)
	})

	copyBtn := ui.CreateCopyButton("Copy", uuidEntry, "UUID copied to clipboard.")

	content := container.NewVBox(
		widget.NewLabel("UUID:"),
		uuidEntry,
		container.NewHBox(generateBtn, copyBtn),
	)

	return container.NewTabItem("UUID", content)
}
