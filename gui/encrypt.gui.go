package gui

import (
	"encoding/base64"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"org.donghyuns.com/secure/keygen/biz"
)

func NewEncryptTab(app fyne.App, window fyne.Window) *container.TabItem {
	ui := NewUIComponents(app, window)
	methods := []string{"AES-CBC", "AES-GCM"}

	methodGroup := ui.CreateRadioGroup(methods, "AES-CBC")
	inputEntry := ui.CreateMultiLineEntry("Enter plaintext to encrypt", 3)
	keyEntry := ui.CreateMultiLineEntry("Enter AES key (Hex/Base64) or generate new key", 2)
	resultEntry := ui.CreateMultiLineEntry("Result will appear here...", 3)

	generateKeyBtn := widget.NewButton("Generate New Key", func() {
		GenerateKeyDialog(app, window, keyEntry)
	})

	processBtn := ui.CreateProcessButton("Encrypt", func() {
		selectedMethod := methodGroup.Selected
		text := inputEntry.Text

		if err := ui.ValidateInput(text, "text"); err != nil {
			ui.ShowError(err)
			return
		}

		switch selectedMethod {
		case "AES-CBC", "AES-GCM":
			key, err := DecodeKey(keyEntry.Text)
			if err != nil {
				ui.ShowError(err)
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
				ui.ShowError(errEncrypt)
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

	copyBtn := ui.CreateCopyButton("Copy", resultEntry, "Result copied to clipboard.")

	content := container.NewVBox(
		widget.NewLabel("Select Method:"),
		methodGroup,
		widget.NewLabel("Input Text:"),
		inputEntry,
		widget.NewLabel("AES Key (Hex/Base64):"),
		container.NewVBox(keyEntry, generateKeyBtn),
		processBtn,
		widget.NewLabel("Result:"),
		resultEntry,
		copyBtn,
	)

	return container.NewTabItem("Encryption", content)
}
