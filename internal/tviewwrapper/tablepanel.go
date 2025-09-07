package tviewwrapper

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/stalwartgiraffe/cmr/events"
)

type TablePanel struct {
	*tview.Table
	style *Style

	onCellSelected events.Event[CellParams]
}

func NewTablePanel(ptc tview.TableContent, stop StopFunc, style *Style) *TablePanel {
	p := &TablePanel{
		Table: tview.NewTable(),
		style: style,
	}
	p.Table.SetContent(ptc)
	p.setupTableLayout(stop)
	p.setupEvents()
    p.SetBlurred()
	return p
}

func (p *TablePanel) setupTableLayout(stop StopFunc) {
	p.Select(0, 0).
		SetFixed(1, 1).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				stop()
			}
			if key == tcell.KeyEnter {
				p.SetSelectable(true, true)
			}
		}).SetSelectedFunc(func(row int, column int) {
		p.GetCell(row, column).SetTextColor(tcell.ColorRed)
		p.SetSelectable(true, true)
	})
}

func (p *TablePanel) setupEvents() {
	p.SetSelectedFunc(func(row, col int) {
		p.onCellSelected.Notify(CellParams{row, col})
	})
}

type CellParams struct {
	Row, Col int
}

func (p *TablePanel) OnCellSelectedSubscribe(fn func(CellParams)) {
	p.onCellSelected.Subscribe(fn)
}
func (p *TablePanel) OnCellSelectedNotify(c CellParams) {
	p.onCellSelected.Notify(c)
}

func (p *TablePanel) SetBlurred() {
	p.SetBackgroundColor(p.style.BlurBackground)
}

func (p *TablePanel) SetFocus(tviewApp *tview.Application) {
	p.SetBackgroundColor(p.style.FocusBackground)
	tviewApp.SetFocus(p)
}
