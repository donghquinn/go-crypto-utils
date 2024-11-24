package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"org.donghyuns.com/secure/keygen/biz"
)

// NewRandomStringTab creates a tab for generating random strings with a slider.
func NewRandomStringTab(app fyne.App, window fyne.Window) *container.TabItem {
	lengthLabel := widget.NewLabel("Length: 16 bytes")
	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder("Random string will appear here...")
	//resultEntry.SetReadOnly(true)
	resultEntry.SetMinRowsVisible(2)

	// Slider for selecting byte length
	slider := widget.NewSlider(1, 128)
	slider.SetValue(16) // Default to 16 bytes

	// Generate random string dynamically on slider change
	slider.OnChanged = func(value float64) {
		length := int(value)
		lengthLabel.SetText(fmt.Sprintf("Length: %d bytes", length))

		randomString, err := biz.GenerateCustomRandomString(length)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error generating random string: %v", err), window)
			return
		}
		resultEntry.SetText(randomString)
	}

	// Copy Button
	copyBtn := widget.NewButton("Copy Random String", func() {
		if resultEntry.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", window)
			return
		}
		app.Driver().AllWindows()[0].Clipboard().SetContent(resultEntry.Text)
		dialog.ShowInformation("Copied", "Random string copied to clipboard.", window)
	})

	content := container.NewVBox(
		widget.NewLabel("Random String Generator"),
		lengthLabel,
		slider,
		resultEntry,
		copyBtn,
	)

	return container.NewTabItem("Random String", content)
}
