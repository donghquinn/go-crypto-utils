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

// GenerateKeyDialog displays a dialog to generate an AES key in Hex or Base64 format
func GenerateKeyDialog(app fyne.App, window fyne.Window, keyEntry *widget.Entry) {
	// Define key formats
	keyFormats := []string{"Hexadecimal", "Base64"}
	choiceGroup := widget.NewRadioGroup(keyFormats, nil)
	choiceGroup.SetSelected("Hexadecimal") // Set a default selection to prevent nil

	// Create dialog content
	generateButton := widget.NewButton("Generate", nil)
	cancelButton := widget.NewButton("Cancel", nil)

	// Placeholder for the dialog instance
	var dlg dialog.Dialog

	// Generate button logic
	generateButton.OnTapped = func() {
		selectedFormat := choiceGroup.Selected
		if selectedFormat == "" {
			dialog.ShowError(fmt.Errorf("Please select a key format"), window)
			return
		}

		keyLength := 32 // Default to AES-256; adjust if needed
		key, err := biz.GenerateRandomAESKey(keyLength)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error generating AES key: %v", err), window)
			return
		}

		var keyStr string
		if selectedFormat == "Hexadecimal" {
			keyStr = hex.EncodeToString(key)
		} else if selectedFormat == "Base64" {
			keyStr = base64.StdEncoding.EncodeToString(key)
		}

		keyEntry.SetText(keyStr)
		if dlg != nil {
			dlg.Hide() // Close the dialog
		}
	}

	// Cancel button logic
	cancelButton.OnTapped = func() {
		if dlg != nil {
			dlg.Hide() // Close the dialog
		}
	}

	// Layout for the dialog
	dialogContent := container.NewVBox(
		widget.NewLabel("Select Key Format:"),
		choiceGroup,
		container.NewHBox(generateButton, cancelButton),
	)

	// Create and show the custom dialog
	dlg = dialog.NewCustom("Generate AES Key", "Close", dialogContent, window)
	dlg.Show()
}

func InputKeyDialog(app fyne.App, window fyne.Window, keyEntry *widget.Entry) {
	// Define key formats
	keyFormats := []string{"Hexadecimal", "Base64"}
	choiceGroup := widget.NewRadioGroup(keyFormats, nil)
	choiceGroup.SetSelected("Hexadecimal") // Set a default selection to prevent nil

	// Create dialog content
	// generateButton := widget.NewButton("Generate", nil)
	// cancelButton := widget.NewButton("Cancel", nil)

	// Placeholder for the dialog instance
	var dlg dialog.Dialog

	// Generate button logic
	// generateButton.OnTapped = func() {
	// 	selectedFormat := choiceGroup.Selected
	// 	if selectedFormat == "" {
	// 		dialog.ShowError(fmt.Errorf("Please select a key format"), window)
	// 		return
	// 	}

	// 	// keyLength := 32 // Default to AES-256; adjust if needed
	// 	// key, err := biz.GenerateRandomAESKey(keyLength)
	// 	// if err != nil {
	// 	// 	dialog.ShowError(fmt.Errorf("Error generating AES key: %v", err), window)
	// 	// 	return
	// 	// }

	// 	var keyStr string
	// 	// if selectedFormat == "Hexadecimal" {
	// 	// 	keyStr = hex.EncodeToString(key)
	// 	// } else if selectedFormat == "Base64" {
	// 	// 	keyStr = base64.StdEncoding.EncodeToString(key)
	// 	// }

	// 	keyEntry.SetText(keyStr)
	// 	if dlg != nil {
	// 		dlg.Hide() // Close the dialog
	// 	}
	// }

	// // Cancel button logic
	// cancelButton.OnTapped = func() {
	// 	if dlg != nil {
	// 		dlg.Hide() // Close the dialog
	// 	}
	// }

	// Layout for the dialog
	dialogContent := container.NewVBox(
		widget.NewLabel("Select Key Format:"),
		choiceGroup,
		// container.NewHBox(generateButton, cancelButton),
	)

	// Create and show the custom dialog
	dlg = dialog.NewCustom("Generate AES Key", "Close", dialogContent, window)
	dlg.Show()
}

// decodeKey decodes the key from Hex or Base64 format
func DecodeKey(keyStr string) ([]byte, error) {
	// Try Hex decoding first
	key, err := hex.DecodeString(keyStr)
	if err == nil {
		return key, nil
	}

	// Try Base64 decoding
	key, err = base64.StdEncoding.DecodeString(keyStr)
	if err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("key must be in Hexadecimal or Base64 format")
}
