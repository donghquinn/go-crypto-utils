package gui

import (
	"encoding/base64"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"org.donghyuns.com/secure/keygen/biz"
)

// NewDecryptTab creates the Decryption tab
func NewDecryptTab(app fyne.App, window fyne.Window) *container.TabItem {
	methods := []string{"AES-CBC", "AES-GCM"}
	methodGroup := widget.NewRadioGroup(methods, nil)
	methodGroup.SetSelected("AES-CBC")

	encryptedEntry := widget.NewEntry()
	encryptedEntry.SetPlaceHolder("Enter encrypted data (Base64)")

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Enter AES key (Hex/Base64)")

	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("Decrypted text will appear here...")

	decryptBtn := widget.NewButton("Decrypt", func() {
		selectedMethod := methodGroup.Selected
		encryptedBase64 := encryptedEntry.Text
		keyInput := keyEntry.Text

		if encryptedBase64 == "" {
			dialog.ShowError(fmt.Errorf("Please enter encrypted data."), window)
			return
		}

		if keyInput == "" {
			dialog.ShowError(fmt.Errorf("Please enter a valid AES key."), window)
			return
		}

		key, err := DecodeKey(keyInput)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		encryptedData, err := base64.StdEncoding.DecodeString(encryptedBase64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Invalid Base64 data."), window)
			return
		}

		var decrypted []byte
		var errDecrypt error
		if selectedMethod == "AES-CBC" {
			decrypted, errDecrypt = biz.DecryptAES256CBC(string(encryptedData), key)
		} else {
			decrypted, errDecrypt = biz.DecryptAES256GCM(encryptedData, key)
		}

		if errDecrypt != nil {
			dialog.ShowError(errDecrypt, window)
			return
		}

		resultEntry.SetText(string(decrypted))
	})

	copyBtn := widget.NewButton("Copy", func() {
		if resultEntry.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(resultEntry.Text)
		dialog.ShowInformation("Copied", "Decrypted text copied to clipboard.", window)
	})

	content := container.NewVBox(
		widget.NewLabel("Select Method:"),
		methodGroup,
		widget.NewLabel("Encrypted Data (Base64):"),
		encryptedEntry,
		widget.NewLabel("AES Key (Hex/Base64):"),
		keyEntry,
		decryptBtn,
		widget.NewLabel("Decrypted Text:"),
		resultEntry,
		copyBtn,
	)

	return container.NewTabItem("Decryption", content)
}
