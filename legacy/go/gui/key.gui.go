package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"org.donghyuns.com/secure/keygen/biz"
)

// NewKeyGenTab creates the Key Generation tab
func NewKeyGenTab(app fyne.App, window fyne.Window) *container.TabItem {
	ui := NewUIComponents(app, window)
	keyTypes := []string{"AES-128 (16bytes)", "AES-192 (24bytes)", "AES-256 (32bytes)"}
	keyTypeGroup := ui.CreateRadioGroup(keyTypes, "AES-256 (32bytes)")

	hexEntry := ui.CreateReadOnlyMultiLineEntry("AES Key in Hexadecimal", 2)
	base64Entry := ui.CreateReadOnlyMultiLineEntry("AES Key in Base64", 2)

	generateBtn := ui.CreateProcessButton("Generate Key", func() {
		selected := keyTypeGroup.Selected
		var length int
		switch selected {
		case "AES-128 (16bytes)":
			length = 16
		case "AES-192 (24bytes)":
			length = 24
		case "AES-256 (32bytes)":
			length = 32
		default:
			length = 32
		}

		key, err := biz.GenerateRandomAESKey(length)
		if err != nil {
			ui.ShowError(fmt.Errorf("Key generation failed: %v", err))
			return
		}

		hexKey, base64Key := biz.GenKey(key)
		hexEntry.SetText(hexKey)
		base64Entry.SetText(base64Key)
	})

	copyHexBtn := ui.CreateCopyButton("Copy Hex Key", hexEntry, "Hex key copied to clipboard.")

	copyBase64Btn := ui.CreateCopyButton("Copy Base64 Key", base64Entry, "Base64 key copied to clipboard.")

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
