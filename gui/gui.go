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

// EncryptionMethod represents the available encryption/hash methods
type EncryptionMethod string

const (
	AES EncryptionMethod = "AES"
	SHA EncryptionMethod = "SHA"
)

// BuildGUI constructs the GUI components and returns the main container
func BuildGUI(app fyne.App) fyne.CanvasObject {
	// Create radio buttons for selecting the encryption method
	methodGroup := widget.NewRadioGroup([]string{string(AES), string(SHA)}, nil)
	methodGroup.SetSelected(string(AES)) // Default selection

	// Create an input field for user to enter text
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Enter text to encrypt/hash")

	// Create a button to generate the encrypted/hashed value
	generateButton := widget.NewButton("Generate", nil)

	// Create a read-only entry to display the result
	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("Encrypted/Hashed value will appear here...")

	// Create a button to copy the result to the clipboard
	copyButton := widget.NewButton("Copy", nil)

	// Define the action for the Generate button
	generateButton.OnTapped = func() {
		selectedMethod := EncryptionMethod(methodGroup.Selected)
		text := inputEntry.Text

		if text == "" {
			dialog.ShowError(fmt.Errorf("please enter text to encrypt/hash"), app.Driver().AllWindows()[0])
			return
		}

		switch selectedMethod {
		case AES:
			// Generate AES key
			newKey, err := biz.GenerateRandomAES256Key(32)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Error generating AES key: %v", err), app.Driver().AllWindows()[0])
				return
			}

			// Encrypt the input text
			plaintextBytes := []byte(text)

			encryptedBytes, err := biz.EncryptAES256CBC(plaintextBytes, newKey)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Encryption Error: %v", err), app.Driver().AllWindows()[0])
				return
			}

			// Encode the encrypted bytes to Base64 for display
			encryptedBase64 := base64.StdEncoding.EncodeToString([]byte(encryptedBytes))
			resultEntry.SetText(encryptedBase64)

		case SHA:
			// Hash the input text using SHA256
			hashedBytes, err := biz.HashSHA256([]byte(text))
			if err != nil {
				dialog.ShowError(fmt.Errorf("Hashing Error: %v", err), app.Driver().AllWindows()[0])
				return
			}

			// Encode the hashed bytes to hexadecimal for display
			hashedHex := fmt.Sprintf("%x", hashedBytes)
			resultEntry.SetText(hashedHex)

		default:
			dialog.ShowError(fmt.Errorf("selected method is not supported"), app.Driver().AllWindows()[0])
		}
	}

	// Define the action for the Copy button
	copyButton.OnTapped = func() {
		if resultEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("there is no content to copy"), app.Driver().AllWindows()[0])
			return
		}
		app.SendNotification(&fyne.Notification{
			Title:   "Copied",
			Content: "The encrypted/hashed value has been copied to the clipboard.",
		})
		app.Driver().AllWindows()[0].Clipboard().SetContent(resultEntry.Text)
	}

	// Layout the components
	content := container.NewVBox(
		widget.NewLabel("Select Encryption Method:"),
		methodGroup,
		widget.NewLabel("Input Text:"),
		inputEntry,
		generateButton,
		widget.NewLabel("Result:"),
		resultEntry,
		copyButton,
	)

	return content
}

// showError displays an error dialog
func showError(app fyne.App, title string, err error) {
	dialog := widget.NewLabel(fmt.Sprintf("%s: %v", title, err))
	errorWindow := app.NewWindow(title)
	errorWindow.SetContent(container.NewVBox(
		dialog,
		widget.NewButton("Close", func() {
			errorWindow.Close()
		}),
	))
	errorWindow.Resize(fyne.NewSize(300, 100))
	errorWindow.Show()
}
