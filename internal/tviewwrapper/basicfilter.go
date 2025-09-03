package tviewwrapper

import (
	"github.com/rivo/tview"
)

// BasicFilter implements FilterPanel with a simple input field
type BasicFilter struct {
	*tview.InputField
	onChanged func(string)
}

// NewBasicFilter creates a basic filter input
func NewBasicFilter(placeholder string) *BasicFilter {
	input := tview.NewInputField()
	input.SetLabel("Filter: ")
	input.SetPlaceholder(placeholder)
	input.SetFieldWidth(0) // Use available width

	filter := &BasicFilter{
		InputField: input,
	}

	input.SetChangedFunc(func(text string) {
		if filter.onChanged != nil {
			filter.onChanged(text)
		}
	})

	return filter
}

func (f *BasicFilter) GetPrimitive() tview.Primitive {
	return f.InputField
}

func (f *BasicFilter) GetFilter() string {
	return f.InputField.GetText()
}

func (f *BasicFilter) SetFilter(text string) {
	f.InputField.SetText(text)
}

func (f *BasicFilter) SetChangedFunc(fn func(string)) {
	f.onChanged = fn
}
