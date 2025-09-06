package tviewwrapper

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/stalwartgiraffe/cmr/events"
)

// BasicFilterPanel implements FilterPanel with a simple input field
type BasicFilterPanel struct {
	*tview.InputField

	OnChange events.Event[string]
}

// NewBasicFilterPanel creates a basic filter input
func NewBasicFilterPanel(placeholder string) *BasicFilterPanel {
	input := tview.NewInputField()
	input.SetLabel("Filter: ")
	input.SetPlaceholder(placeholder)
	input.SetFieldWidth(0) // Use available width

	//input.SetFieldTextColor(tcell.ColorDarkSlateGray)
	//input.SetFieldBackgroundColor(tcell.ColorDarkGrey)
	//almostBlack := tcell.NewRGBColor(10, 10, 10)
	almostBlack := tcell.NewRGBColor(47, 47, 47)
	input.SetFieldBackgroundColor(almostBlack)
	//input.SetFieldBackgroundColor(tcell.ColorDarkGrey)

	f := &BasicFilterPanel{
		InputField: input,
	}

	input.SetChangedFunc(func(text string) {
		f.OnChange.Notify(text)
	})

	return f
}

func (f *BasicFilterPanel) OnChangeSubscribe(fn func(string)) {
	f.OnChange.Subscribe(fn)
}

func (f *BasicFilterPanel) GetPrimitive() tview.Primitive {
	return f.InputField
}

func (f *BasicFilterPanel) GetFilter() string {
	return f.GetText()
}

func (f *BasicFilterPanel) SetFilter(text string) {
	f.SetText(text)
}
