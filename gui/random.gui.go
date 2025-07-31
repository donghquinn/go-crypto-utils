package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"org.donghyuns.com/secure/keygen/biz"
)

// NewRandomStringTab creates a tab for generating random strings with a slider.
func NewRandomStringTab(app fyne.App, window fyne.Window) *container.TabItem {
	ui := NewUIComponents(app, window)
	lengthLabel := widget.NewLabel("Length: 16 bytes")
	resultEntry := ui.CreateMultiLineEntry("Random string will appear here...", 2)

	// Slider for selecting byte length
	slider := widget.NewSlider(1, 128)
	slider.SetValue(16) // Default to 16 bytes

	// Generate random string dynamically on slider change
	slider.OnChanged = func(value float64) {
		length := int(value)
		lengthLabel.SetText(fmt.Sprintf("Length: %d bytes", length))

		randomString, err := biz.GenerateCustomRandomString(length)
		if err != nil {
			ui.ShowError(fmt.Errorf("Error generating random string: %v", err))
			return
		}
		resultEntry.SetText(randomString)
	}

	// Copy Button
	copyBtn := ui.CreateCopyButton("Copy Random String", resultEntry, "Random string copied to clipboard.")

	content := container.NewVBox(
		widget.NewLabel("Random String Generator"),
		lengthLabel,
		slider,
		resultEntry,
		copyBtn,
	)

	return container.NewTabItem("Random String", content)
}
