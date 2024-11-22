package gui

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewEncodeTab(app fyne.App, window fyne.Window) *container.TabItem {
	// Base64 Encode Tab
	base64Input := widget.NewMultiLineEntry()
	base64Input.SetPlaceHolder("Enter text to encode into Base64...")
	base64Input.SetMinRowsVisible(3)

	base64Result := widget.NewMultiLineEntry()
	base64Result.SetPlaceHolder("Base64 encoded result will appear here...")
	// base64Result.SetReadOnly(true)
	base64Result.SetMinRowsVisible(3)

	base64EncodeBtn := widget.NewButton("Encode", func() {
		input := base64Input.Text
		if input == "" {
			dialog.ShowError(fmt.Errorf("Please enter text to encode."), window)
			return
		}

		base64Result.SetText(base64.StdEncoding.EncodeToString([]byte(input)))
	})

	base64CopyBtn := widget.NewButton("Copy Result", func() {
		if base64Result.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(base64Result.Text)
		dialog.ShowInformation("Copied", "Encoded text copied to clipboard.", window)
	})

	base64Tab := container.NewVBox(
		widget.NewLabel("Base64 Encode"),
		base64Input,
		base64EncodeBtn,
		base64Result,
		base64CopyBtn,
	)

	// Hex Encode Tab
	hexInput := widget.NewMultiLineEntry()
	hexInput.SetPlaceHolder("Enter text to encode into Hex...")
	hexInput.SetMinRowsVisible(3)

	hexResult := widget.NewMultiLineEntry()
	hexResult.SetPlaceHolder("Hex encoded result will appear here...")
	// hexResult.SetReadOnly(true)
	hexResult.SetMinRowsVisible(3)

	hexEncodeBtn := widget.NewButton("Encode", func() {
		input := hexInput.Text
		if input == "" {
			dialog.ShowError(fmt.Errorf("Please enter text to encode."), window)
			return
		}

		hexResult.SetText(hex.EncodeToString([]byte(input)))
	})

	hexCopyBtn := widget.NewButton("Copy Result", func() {
		if hexResult.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(hexResult.Text)
		dialog.ShowInformation("Copied", "Encoded text copied to clipboard.", window)
	})

	hexTab := container.NewVBox(
		widget.NewLabel("Hex Encode"),
		hexInput,
		hexEncodeBtn,
		hexResult,
		hexCopyBtn,
	)

	// Main Encode Tabs
	encodeTabs := container.NewAppTabs(
		container.NewTabItem("Base64", base64Tab),
		container.NewTabItem("Hex", hexTab),
	)

	return container.NewTabItem("Encode", encodeTabs)
}

func NewDecodeTab(app fyne.App, window fyne.Window) *container.TabItem {
	// Base64 Decode Tab
	base64Input := widget.NewMultiLineEntry()
	base64Input.SetPlaceHolder("Enter Base64 string to decode...")
	base64Input.SetMinRowsVisible(3)

	base64Result := widget.NewMultiLineEntry()
	base64Result.SetPlaceHolder("Decoded result will appear here...")
	// base64Result.SetReadOnly(true)
	base64Result.SetMinRowsVisible(3)

	base64DecodeBtn := widget.NewButton("Decode", func() {
		input := base64Input.Text
		if input == "" {
			dialog.ShowError(fmt.Errorf("Please enter a Base64 string to decode."), window)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error decoding Base64 string: %v", err), window)
			return
		}

		base64Result.SetText(string(decoded))
	})

	base64CopyBtn := widget.NewButton("Copy Result", func() {
		if base64Result.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(base64Result.Text)
		dialog.ShowInformation("Copied", "Decoded text copied to clipboard.", window)
	})

	base64Tab := container.NewVBox(
		widget.NewLabel("Base64 Decode"),
		base64Input,
		base64DecodeBtn,
		base64Result,
		base64CopyBtn,
	)

	// Hex Decode Tab
	hexInput := widget.NewMultiLineEntry()
	hexInput.SetPlaceHolder("Enter Hex string to decode...")
	hexInput.SetMinRowsVisible(3)

	hexResult := widget.NewMultiLineEntry()
	hexResult.SetPlaceHolder("Decoded result will appear here...")
	// hexResult.SetReadOnly(true)
	hexResult.SetMinRowsVisible(3)

	hexDecodeBtn := widget.NewButton("Decode", func() {
		input := hexInput.Text
		if input == "" {
			dialog.ShowError(fmt.Errorf("Please enter a Hex string to decode."), window)
			return
		}

		decoded, err := hex.DecodeString(input)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error decoding Hex string: %v", err), window)
			return
		}

		hexResult.SetText(string(decoded))
	})

	hexCopyBtn := widget.NewButton("Copy Result", func() {
		if hexResult.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(hexResult.Text)
		dialog.ShowInformation("Copied", "Decoded text copied to clipboard.", window)
	})

	hexTab := container.NewVBox(
		widget.NewLabel("Hex Decode"),
		hexInput,
		hexDecodeBtn,
		hexResult,
		hexCopyBtn,
	)

	// Main Decode Tabs
	decodeTabs := container.NewAppTabs(
		container.NewTabItem("Base64", base64Tab),
		container.NewTabItem("Hex", hexTab),
	)

	return container.NewTabItem("Decode", decodeTabs)
}
