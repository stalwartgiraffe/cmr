package tviewwrapper

import (
	"fmt"
	"reflect"

	"github.com/rivo/tview"
)

// TextDetailsPanel implements DetailsPanel with a text view
type TextDetailsPanel struct {
	*tview.TextView
}

// NewTextDetails creates a basic text details panel
func NewTextDetailsPanel() *TextDetailsPanel {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle("Details")
	textView.SetDynamicColors(true)
	textView.SetWordWrap(true)
	textView.SetScrollable(true)

	p := &TextDetailsPanel{
		TextView: textView,
	}

	return p
}

func (d *TextDetailsPanel) GetPrimitive() tview.Primitive {
	return d.TextView
}

func (d *TextDetailsPanel) ShowDetails(data any) {
	if data == nil {
		d.Clear()
		return
	}

	d.TextView.SetText(formatDetails(data))
}

func (d *TextDetailsPanel) Clear() {
	d.TextView.SetText("")
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
