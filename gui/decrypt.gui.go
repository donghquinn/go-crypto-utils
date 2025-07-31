package gui

import (
	"encoding/base64"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"org.donghyuns.com/secure/keygen/biz"
)

// NewDecryptTab creates the Decryption tab
func NewDecryptTab(app fyne.App, window fyne.Window) *container.TabItem {
	ui := NewUIComponents(app, window)
	methods := []string{"AES-CBC", "AES-GCM"}
	methodGroup := ui.CreateRadioGroup(methods, "AES-CBC")

	encryptedEntry := widget.NewEntry()
	encryptedEntry.SetPlaceHolder("Enter encrypted data (Hex/Base64)")

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Enter key (Hex/Base64)")

	resultEntry := ui.CreateMultiLineEntry("Decrypted text will appear here...", 3)

	decryptBtn := ui.CreateProcessButton("Decrypt", func() {
		selectedMethod := methodGroup.Selected
		encryptedBase64 := encryptedEntry.Text
		keyInput := keyEntry.Text

		if err := ui.ValidateInput(encryptedBase64, "encrypted data"); err != nil {
			ui.ShowError(err)
			return
		}

		if err := ui.ValidateInput(keyInput, "a valid AES key"); err != nil {
			ui.ShowError(err)
			return
		}

		key, err := DecodeKey(keyInput)
		if err != nil {
			ui.ShowError(err)
			return
		}

		encryptedData, err := base64.StdEncoding.DecodeString(encryptedBase64)
		if err != nil {
			ui.ShowError(fmt.Errorf("Invalid Base64 data."))
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
			ui.ShowError(errDecrypt)
			return
		}

		resultEntry.SetText(string(decrypted))
	})

	copyBtn := ui.CreateCopyButton("Copy", resultEntry, "Decrypted text copied to clipboard.")

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
