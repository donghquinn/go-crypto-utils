package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"org.donghyuns.com/secure/keygen/biz"
)

// NewKeyGenTab creates the Key Generation tab
func NewKeyGenTab(app fyne.App, window fyne.Window) *container.TabItem {
	keyTypes := []string{"AES-128 (16bytes)", "AES-192 (24bytes)", "AES-256 (32bytes)"}
	keyTypeGroup := widget.NewRadioGroup(keyTypes, nil)
	keyTypeGroup.SetSelected("AES-256")

	// Read-only multi-line entry for Hex key
	hexEntry := widget.NewMultiLineEntry()
	hexEntry.SetPlaceHolder("AES Key in Hexadecimal")
	hexEntry.Disable()                     // Make the box read-only
	hexEntry.SetMinRowsVisible(2)          // Increase height
	hexEntry.Wrapping = fyne.TextWrapBreak // Enable text wrapping

	// Read-only multi-line entry for Base64 key
	base64Entry := widget.NewMultiLineEntry()
	base64Entry.SetPlaceHolder("AES Key in Base64")
	base64Entry.Disable()                     // Make the box read-only
	base64Entry.SetMinRowsVisible(2)          // Increase height
	base64Entry.Wrapping = fyne.TextWrapBreak // Enable text wrapping

	generateBtn := widget.NewButton("Generate Key", func() {
		selected := keyTypeGroup.Selected
		var length int
		switch selected {
		case "AES-128 (16bytes)":
			length = 16
		case "AES-192 (24bytes)":
			length = 24
		case "AES-256 (32bytes)":
			length = 32
		}

		key, err := biz.GenerateRandomAESKey(length)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Key generation failed: %v", err), window)
			return
		}

		hexKey, base64Key := biz.GenKey(key)
		hexEntry.SetText(hexKey)
		base64Entry.SetText(base64Key)
	})

	copyHexBtn := widget.NewButton("Copy Hex Key", func() {
		if hexEntry.Text == "" {
			dialog.ShowInformation("No Key", "No Hex key to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(hexEntry.Text)
		dialog.ShowInformation("Copied", "Hex key copied to clipboard.", window)
	})

	copyBase64Btn := widget.NewButton("Copy Base64 Key", func() {
		if base64Entry.Text == "" {
			dialog.ShowInformation("No Key", "No Base64 key to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(base64Entry.Text)
		dialog.ShowInformation("Copied", "Base64 key copied to clipboard.", window)
	})

	content := container.NewVBox(
		widget.NewLabel("Select Key Type:"),
		keyTypeGroup,
		generateBtn,
		widget.NewLabel("AES Key (Hexadecimal):"),
		container.NewVBox(hexEntry, copyHexBtn), // Adjusted layout for better readability
		widget.NewLabel("AES Key (Base64):"),
		container.NewVBox(base64Entry, copyBase64Btn), // Adjusted layout for better readability
	)

	return container.NewTabItem("Key Generation", content)
}
