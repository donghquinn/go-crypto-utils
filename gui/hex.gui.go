package gui

import (
	"encoding/hex"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewEncodeHexTab(app fyne.App, window fyne.Window) *container.TabItem {
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Enter text to encode into Hex...")
	inputEntry.SetMinRowsVisible(3)

	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("Encoded Hex result will appear here...")
	// resultEntry.SetReadOnly(true)
	resultEntry.SetMinRowsVisible(3)

	encodeBtn := widget.NewButton("Encode", func() {
		input := inputEntry.Text
		if input == "" {
			dialog.ShowError(fmt.Errorf("Please enter text to encode."), window)
			return
		}

		encoded := hex.EncodeToString([]byte(input))
		resultEntry.SetText(encoded)
	})

	copyBtn := widget.NewButton("Copy Result", func() {
		if resultEntry.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(resultEntry.Text)
		dialog.ShowInformation("Copied", "Encoded text copied to clipboard.", window)
	})

	content := container.NewVBox(
		widget.NewLabel("Encode to Hex"),
		widget.NewLabel("Input:"),
		inputEntry,
		encodeBtn,
		widget.NewLabel("Result:"),
		resultEntry,
		copyBtn,
	)

	return container.NewTabItem("Encode Hex", content)
}

func NewDecodeHexTab(app fyne.App, window fyne.Window) *container.TabItem {
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Enter Hex-encoded string here...")
	inputEntry.SetMinRowsVisible(3)

	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("Decoded result will appear here...")
	// resultEntry.SetReadOnly(true)
	resultEntry.SetMinRowsVisible(3)

	decodeBtn := widget.NewButton("Decode", func() {
		input := inputEntry.Text
		if input == "" {
			dialog.ShowError(fmt.Errorf("Please enter a Hex string to decode."), window)
			return
		}

		decoded, err := hex.DecodeString(input)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error decoding Hex string: %v", err), window)
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
		dialog.ShowInformation("Copied", "Decoded text copied to clipboard.", window)
	})

	content := container.NewVBox(
		widget.NewLabel("Decode from Hex"),
		widget.NewLabel("Input:"),
		inputEntry,
		decodeBtn,
		widget.NewLabel("Result:"),
		resultEntry,
		copyBtn,
	)

	return container.NewTabItem("Decode Hex", content)
}
