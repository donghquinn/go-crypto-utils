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
	keyTypes := []string{"AES-128", "AES-192", "AES-256"}
	keyTypeGroup := widget.NewRadioGroup(keyTypes, nil)
	keyTypeGroup.SetSelected("AES-256")

	hexEntry := widget.NewEntry()
	hexEntry.SetPlaceHolder("AES Key in Hexadecimal")

	base64Entry := widget.NewEntry()
	base64Entry.SetPlaceHolder("AES Key in Base64")

	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("Result will appear here...")

	generateBtn := widget.NewButton("Generate Key", func() {
		selected := keyTypeGroup.Selected
		var length int
		switch selected {
		case "AES-128":
			length = 16
		case "AES-192":
			length = 24
		case "AES-256":
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
		app.Driver().AllWindows()[0].Clipboard().SetContent(resultEntry.Text)
		dialog.ShowInformation("Copied", "Hex key copied to clipboard.", window)
	})

	copyBase64Btn := widget.NewButton("Copy Base64 Key", func() {
		if base64Entry.Text == "" {
			dialog.ShowInformation("No Key", "No Base64 key to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(resultEntry.Text)
		dialog.ShowInformation("Copied", "Base64 key copied to clipboard.", window)
	})

	content := container.NewVBox(
		widget.NewLabel("Select Key Type:"),
		keyTypeGroup,
		generateBtn,
		widget.NewLabel("AES Key (Hexadecimal):"),
		container.NewHBox(hexEntry, copyHexBtn),
		widget.NewLabel("AES Key (Base64):"),
		container.NewHBox(base64Entry, copyBase64Btn),
	)

	return container.NewTabItem("Key Generation", content)
}
