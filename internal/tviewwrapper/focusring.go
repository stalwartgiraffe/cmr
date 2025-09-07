package tviewwrapper

import (
	"github.com/rivo/tview"

	"github.com/stalwartgiraffe/cmr/events"
)

// FocusRing implements cycling focus through the slice of panels in a ring.
type FocusRing struct {
	tviewApp     *tview.Application
	focusedPanel int
	panels       []Panel

	onPanelBlurred events.Event[FocusParams]
	onPanelFocused events.Event[FocusParams]
}

type Panel interface {
	tview.Primitive
	Focusable
}

type Focusable interface {
	SetBlurred()
	SetFocus(tviewApp *tview.Application)
}

// NewFocusRing returns a configured focus ring of panels.
// func NewFocusRing(tviewApp *tview.Application, panels ...tview.Primitive) *FocusRing {
func NewFocusRing(tviewApp *tview.Application, panels ...Panel) *FocusRing {
	r := &FocusRing{
		tviewApp:     tviewApp,
		focusedPanel: 0,
		panels:       panels,
	}
	r.onPanelBlurred.Subscribe(func(p FocusParams) {
		p.Panel.SetBlurred()
	})
	r.onPanelFocused.Subscribe(func(p FocusParams) {
		p.Panel.SetFocus(p.TviewApp)
	})
	return r
}

type FocusParams struct {
	TviewApp *tview.Application
	Panel    Panel
}

type RingDirection int

const (
	PrevDir RingDirection = -1
	NextDir RingDirection = 1
)

// Cycle changes the focus in the specified direction
func (r *FocusRing) Cycle(direction RingDirection) {
	N := len(r.panels)

	r.onPanelBlurred.Notify(FocusParams{
		r.tviewApp,
		r.panels[r.focusedPanel],
	})
	r.focusedPanel = ((r.focusedPanel+int(direction))%N + N) % N
	r.onPanelFocused.Notify(FocusParams{
		r.tviewApp,
		r.panels[r.focusedPanel],
	})
}
