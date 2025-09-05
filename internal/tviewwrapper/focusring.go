package tviewwrapper

import (
	"github.com/rivo/tview"

	"github.com/stalwartgiraffe/cmr/events"
)

// FocusRing implements cycling focus through the slice of panels in a ring.
type FocusRing struct {
	tviewApp     *tview.Application
	focusedPanel int
	panels       []tview.Primitive

	onPanelFocused events.Event[FocusParams]
}

// NewFocusRing returns a configured focus ring of panels.
func NewFocusRing(tviewApp *tview.Application, panels ...tview.Primitive) *FocusRing {
	r := &FocusRing{
		tviewApp:     tviewApp,
		focusedPanel: 0,
		panels:       panels,
	}
	r.onPanelFocused.Subscribe(func(p FocusParams) {
		p.TviewApp.SetFocus(p.Panel)
	})
	return r
}

type FocusParams struct {
	TviewApp *tview.Application
	Panel    tview.Primitive
}

type RingDirection int

const (
	PrevDir RingDirection = -1
	NextDir RingDirection = 1
)

// Cycle changes the focus in the specified direction
func (r *FocusRing) Cycle(direction RingDirection) {
	N := len(r.panels)
	r.focusedPanel = ((r.focusedPanel+int(direction))%N + N) % N
	r.onPanelFocused.Notify(FocusParams{
		r.tviewApp,
		r.panels[r.focusedPanel],
	})
}
