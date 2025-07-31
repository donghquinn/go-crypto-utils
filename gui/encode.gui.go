package gui

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewEncodeTab(app fyne.App, window fyne.Window) *container.TabItem {
	ui := NewUIComponents(app, window)
	
	base64Input := ui.CreateMultiLineEntry("Enter text to encode into Base64...", 3)
	base64Result := ui.CreateMultiLineEntry("Base64 encoded result will appear here...", 3)

	base64EncodeBtn := ui.CreateProcessButton("Encode", func() {
		input := base64Input.Text
		if err := ui.ValidateInput(input, "text to encode"); err != nil {
			ui.ShowError(err)
			return
		}

		base64Result.SetText(base64.StdEncoding.EncodeToString([]byte(input)))
	})

	base64CopyBtn := ui.CreateCopyButton("Copy Result", base64Result, "Encoded text copied to clipboard.")

	base64Tab := container.NewVBox(
		widget.NewLabel("Base64 Encode"),
		base64Input,
		base64EncodeBtn,
		base64Result,
		base64CopyBtn,
	)

	hexInput := ui.CreateMultiLineEntry("Enter text to encode into Hex...", 3)
	hexResult := ui.CreateMultiLineEntry("Hex encoded result will appear here...", 3)

	hexEncodeBtn := ui.CreateProcessButton("Encode", func() {
		input := hexInput.Text
		if err := ui.ValidateInput(input, "text to encode"); err != nil {
			ui.ShowError(err)
			return
		}

		hexResult.SetText(hex.EncodeToString([]byte(input)))
	})

	hexCopyBtn := ui.CreateCopyButton("Copy Result", hexResult, "Encoded text copied to clipboard.")

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
	ui := NewUIComponents(app, window)
	
	base64Input := ui.CreateMultiLineEntry("Enter Base64 string to decode...", 3)
	base64Result := ui.CreateMultiLineEntry("Decoded result will appear here...", 3)

	base64DecodeBtn := ui.CreateProcessButton("Decode", func() {
		input := base64Input.Text
		if err := ui.ValidateInput(input, "a Base64 string to decode"); err != nil {
			ui.ShowError(err)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			ui.ShowError(fmt.Errorf("Error decoding Base64 string: %v", err))
			return
		}

		base64Result.SetText(string(decoded))
	})

	base64CopyBtn := ui.CreateCopyButton("Copy Result", base64Result, "Decoded text copied to clipboard.")

	base64Tab := container.NewVBox(
		widget.NewLabel("Base64 Decode"),
		base64Input,
		base64DecodeBtn,
		base64Result,
		base64CopyBtn,
	)

	hexInput := ui.CreateMultiLineEntry("Enter Hex string to decode...", 3)
	hexResult := ui.CreateMultiLineEntry("Decoded result will appear here...", 3)

	hexDecodeBtn := ui.CreateProcessButton("Decode", func() {
		input := hexInput.Text
		if err := ui.ValidateInput(input, "a Hex string to decode"); err != nil {
			ui.ShowError(err)
			return
		}

		decoded, err := hex.DecodeString(input)
		if err != nil {
			ui.ShowError(fmt.Errorf("Error decoding Hex string: %v", err))
			return
		}

		hexResult.SetText(string(decoded))
	})

	hexCopyBtn := ui.CreateCopyButton("Copy Result", hexResult, "Decoded text copied to clipboard.")

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
