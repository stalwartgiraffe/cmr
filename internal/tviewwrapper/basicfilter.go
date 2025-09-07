package tviewwrapper

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/stalwartgiraffe/cmr/events"
)

// BasicFilterPanel implements FilterPanel with a simple input field
type BasicFilterPanel struct {
	*tview.InputField

	style *Style

	OnChange events.Event[string]
}

// NewBasicFilterPanel creates a basic filter input
func NewBasicFilterPanel(placeholder string, style *Style) *BasicFilterPanel {
	input := tview.NewInputField()
	input.SetLabel("Filter: ")
	input.SetPlaceholder(placeholder)
	input.SetFieldWidth(0) // Use available width

	inputDarkGray := tcell.Color236
	input.SetFieldBackgroundColor(inputDarkGray)

	p := &BasicFilterPanel{
		InputField: input,
		style:      style,
	}
	p.SetBlurred()

	input.SetChangedFunc(func(text string) {
		p.OnChange.Notify(text)
	})

	return p
}

func (p *BasicFilterPanel) OnChangeSubscribe(fn func(string)) {
	p.OnChange.Subscribe(fn)
}

func (p *BasicFilterPanel) GetPrimitive() tview.Primitive {
	return p.InputField
}

func (p *BasicFilterPanel) GetFilter() string {
	return p.GetText()
}

func (p *BasicFilterPanel) SetFilter(text string) {
	p.SetText(text)
}

func (p *BasicFilterPanel) SetBlurred() {
	p.SetBackgroundColor(p.style.BlurBackground)
}

func (p *BasicFilterPanel) SetFocus(tviewApp *tview.Application) {
	p.SetBackgroundColor(p.style.FocusBackground)
	tviewApp.SetFocus(p)
}
