package gui

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"org.donghyuns.com/secure/keygen/biz"
)

// NewEncryptTab creates the Encryption tab
func NewEncryptTab(app fyne.App, window fyne.Window) *container.TabItem {
	methods := []string{"AES-CBC", "AES-GCM", "SHA-256", "SHA-512"}
	methodGroup := widget.NewRadioGroup(methods, nil)
	methodGroup.SetSelected("AES-CBC")

	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Enter text to encrypt/hash")

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Enter AES key (Hex/Base64)")

	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("Result will appear here...")

	generateKeyBtn := widget.NewButton("Generate Key", func() {
		GenerateKeyDialog(app, window, keyEntry)
	})

	processBtn := widget.NewButton("Process", func() {
		selectedMethod := methodGroup.Selected
		text := inputEntry.Text

		if text == "" {
			dialog.ShowError(fmt.Errorf("Please enter text."), window)
			return
		}

		switch selectedMethod {
		case "AES-CBC", "AES-GCM":
			key, err := DecodeKey(keyEntry.Text)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			var encrypted []byte
			var errEncrypt error
			if selectedMethod == "AES-CBC" {
				encryptedStr, _ := biz.EncryptAES256CBC([]byte(text), key)
				encrypted = []byte(encryptedStr)
			} else {
				encrypted, errEncrypt = biz.EncryptAES256GCM([]byte(text), key)
			}

			if errEncrypt != nil {
				dialog.ShowError(errEncrypt, window)
				return
			}

			result := base64.StdEncoding.EncodeToString(encrypted)
			resultEntry.SetText(result)

		case "SHA-256", "SHA-512":
			hashed, err := biz.HashData([]byte(text), selectedMethod)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			result := hex.EncodeToString(hashed)
			resultEntry.SetText(result)
		}
	})

	copyBtn := widget.NewButton("Copy", func() {
		if resultEntry.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(resultEntry.Text)
		dialog.ShowInformation("Copied", "Result copied to clipboard.", window)
	})

	content := container.NewVBox(
		widget.NewLabel("Select Method:"),
		methodGroup,
		widget.NewLabel("Input Text:"),
		inputEntry,
		container.NewHBox(keyEntry, generateKeyBtn),
		processBtn,
		widget.NewLabel("Result:"),
		resultEntry,
		copyBtn,
	)

	return container.NewTabItem("Encryption", content)
}
