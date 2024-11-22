package gui

import (
	"encoding/base64"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// NewDecodeBase64Tab creates the Decode Base64 Strings tab
func NewDecodeBase64Tab(app fyne.App, window fyne.Window) *container.TabItem {
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Enter Base64-encoded string here...")
	inputEntry.SetMinRowsVisible(2)

	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("Decoded result will appear here...")
	// resultEntry.Disable()
	resultEntry.SetMinRowsVisible(2)

	decodeBtn := widget.NewButton("Decode", func() {
		input := inputEntry.Text
		if input == "" {
			dialog.ShowError(fmt.Errorf("Please enter a Base64-encoded string."), window)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error decoding Base64 string: %v", err), window)
			return
		}

		resultEntry.SetText(string(decoded))
	})

	copyBtn := widget.NewButton("Copy Decoded Text", func() {
		if resultEntry.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(resultEntry.Text)
		dialog.ShowInformation("Copied", "Decoded text copied to clipboard.", window)
	})

	content := container.NewVBox(
		widget.NewLabel("Base64 Decoder"),
		widget.NewLabel("Input:"),
		inputEntry,
		decodeBtn,
		widget.NewLabel("Result:"),
		resultEntry,
		copyBtn,
	)

	return container.NewTabItem("Decode Base64", content)
}

// NewBase64Tab creates the Base64 Encoder/Decoder tab
func NewBase64Tab(app fyne.App, window fyne.Window) *container.TabItem {
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Enter text or Base64 string here...")
	inputEntry.SetMinRowsVisible(3)

	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("Result will appear here...")
	// resultEntry.SetReadOnly(true)
	resultEntry.SetMinRowsVisible(3)

	encodeBtn := widget.NewButton("Encode to Base64", func() {
		input := inputEntry.Text
		if input == "" {
			dialog.ShowError(fmt.Errorf("Please enter text to encode."), window)
			return
		}

		encoded := base64.StdEncoding.EncodeToString([]byte(input))
		resultEntry.SetText(encoded)
	})

	decodeBtn := widget.NewButton("Decode from Base64", func() {
		input := inputEntry.Text
		if input == "" {
			dialog.ShowError(fmt.Errorf("Please enter a Base64 string to decode."), window)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error decoding Base64 string: %v", err), window)
			return
		}

		resultEntry.SetText(string(decoded))
	})

	copyBtn := widget.NewButton("Copy Result", func() {
		if resultEntry.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(resultEntry.Text)
		dialog.ShowInformation("Copied", "Result copied to clipboard.", window)
	})

	content := container.NewVBox(
		widget.NewLabel("Base64 Encoder/Decoder"),
		widget.NewLabel("Input:"),
		inputEntry,
		container.NewHBox(encodeBtn, decodeBtn), // Encode and Decode buttons side by side
		widget.NewLabel("Result:"),
		resultEntry,
		copyBtn,
	)

	return container.NewTabItem("Base64 Encoder/Decoder", content)
}
