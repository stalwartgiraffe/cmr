package tviewwrapper

import (
	"fmt"
	"reflect"

	"github.com/rivo/tview"
)

// TextDetailsPanel implements DetailsPanel with a text view
type TextDetailsPanel struct {
	*tview.TextView
	style *Style
}

// NewTextDetails creates a basic text details panel
func NewTextDetailsPanel(style *Style) *TextDetailsPanel {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle("Details")
	textView.SetDynamicColors(true)
	textView.SetWordWrap(true)
	textView.SetScrollable(true)

	p := &TextDetailsPanel{
		TextView: textView,
		style:    style,
	}
	p.SetBlurred()

	return p
}

func (p *TextDetailsPanel) GetPrimitive() tview.Primitive {
	return p.TextView
}

func (p *TextDetailsPanel) ShowDetails(data any) {
	if data == nil {
		p.Clear()
		return
	}

	p.TextView.SetText(formatDetails(data))
}

func (p *TextDetailsPanel) Clear() {
	p.TextView.SetText("")
}

// formatDetails uses reflection to format any struct for display
func formatDetails(data any) string {
	if data == nil {
		return "No details available"
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Sprintf("Value: %v", data)
	}

	t := v.Type()
	result := fmt.Sprintf("[yellow]%s Details[white]\n\n", t.Name())

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if !field.IsExported() {
			continue
		}

		result += fmt.Sprintf("[blue]%s:[white] %v\n", field.Name, value.Interface())
	}

	return result
}

func (p *TextDetailsPanel) SetBlurred() {
	p.SetBackgroundColor(p.style.BlurBackground)
}

func (p *TextDetailsPanel) SetFocus(tviewApp *tview.Application) {
	p.SetBackgroundColor(p.style.FocusBackground)
	tviewApp.SetFocus(p)
}
