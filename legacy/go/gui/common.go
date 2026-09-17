package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	MinRowsForEntry    = 2
	MinRowsForMultiple = 3
	DefaultKeySize     = 32
)

type UIComponents struct {
	App    fyne.App
	Window fyne.Window
}

func NewUIComponents(app fyne.App, window fyne.Window) *UIComponents {
	return &UIComponents{
		App:    app,
		Window: window,
	}
}

func (ui *UIComponents) CreateMultiLineEntry(placeholder string, minRows int) *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder(placeholder)
	entry.SetMinRowsVisible(minRows)
	entry.Wrapping = fyne.TextWrapBreak
	return entry
}

func (ui *UIComponents) CreateReadOnlyMultiLineEntry(placeholder string, minRows int) *widget.Entry {
	entry := ui.CreateMultiLineEntry(placeholder, minRows)
	entry.Disable()
	return entry
}

func (ui *UIComponents) CreateRadioGroup(options []string, defaultSelection string) *widget.RadioGroup {
	group := widget.NewRadioGroup(options, nil)
	if defaultSelection != "" {
		group.SetSelected(defaultSelection)
	}
	return group
}

func (ui *UIComponents) CreateCopyButton(text string, sourceEntry *widget.Entry, successMessage string) *widget.Button {
	return widget.NewButton(text, func() {
		if sourceEntry.Text == "" {
			dialog.ShowInformation("No Content", "Nothing to copy.", ui.Window)
			return
		}
		ui.App.Driver().AllWindows()[0].Clipboard().SetContent(sourceEntry.Text)
		dialog.ShowInformation("Copied", successMessage, ui.Window)
	})
}

func (ui *UIComponents) ShowError(err error) {
	dialog.ShowError(err, ui.Window)
}

func (ui *UIComponents) ShowInfo(title, message string) {
	dialog.ShowInformation(title, message, ui.Window)
}

func (ui *UIComponents) CreateLabeledInput(label string, entry *widget.Entry) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabel(label),
		entry,
	)
}

func (ui *UIComponents) CreateButtonRow(buttons ...*widget.Button) fyne.CanvasObject {
	var items []fyne.CanvasObject
	for _, btn := range buttons {
		items = append(items, btn)
	}
	return container.NewHBox(items...)
}

func (ui *UIComponents) CreateProcessButton(text string, callback func()) *widget.Button {
	return widget.NewButton(text, callback)
}

func (ui *UIComponents) ValidateInput(input, fieldName string) error {
	if input == "" {
		return &ValidationError{Field: fieldName}
	}
	return nil
}

type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return "Please enter " + e.Field + "."
}