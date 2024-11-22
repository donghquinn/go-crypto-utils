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

func NewEncryptTab(app fyne.App, window fyne.Window) *container.TabItem {
	methods := []string{"AES-CBC", "AES-GCM"} // "SHA-256", "SHA-512"

	methodGroup := widget.NewRadioGroup(methods, nil)
	methodGroup.SetSelected("AES-CBC")

	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Enter plaintext to encrypt")

	// Use MultiLineEntry for AES key input to support long keys
	keyEntry := widget.NewMultiLineEntry()
	keyEntry.SetPlaceHolder("Enter AES key (Hex/Base64) or generate new key")
	keyEntry.Wrapping = fyne.TextWrapBreak // Ensure text wraps within the box
	keyEntry.SetMinRowsVisible(2)          // Make the box taller

	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("Result will appear here...")

	generateKeyBtn := widget.NewButton("Generate New Key", func() {
		GenerateKeyDialog(app, window, keyEntry)
	})

	processBtn := widget.NewButton("Encrypt", func() {
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

			// case "SHA-256", "SHA-512":
			// 	hashed, err := biz.HashData([]byte(text), selectedMethod)
			// 	if err != nil {
			// 		dialog.ShowError(err, window)
			// 		return
			// 	}
			// 	result := hex.EncodeToString(hashed)
			// 	resultEntry.SetText(result)
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
		widget.NewLabel("AES Key (Hex/Base64):"),
		container.NewVBox(keyEntry, generateKeyBtn), // Adjust layout
		processBtn,
		widget.NewLabel("Result:"),
		resultEntry,
		copyBtn,
	)

	return container.NewTabItem("Encryption", content)
}
