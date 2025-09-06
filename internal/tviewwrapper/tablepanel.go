package tviewwrapper

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/stalwartgiraffe/cmr/events"
)

type TablePanel struct {
	*tview.Table

	onCellSelected events.Event[CellParams]
}

func NewTablePanel(ptc tview.TableContent, stop StopFunc) *TablePanel {
	t := &TablePanel{
		Table: tview.NewTable(),
	}
	t.Table.SetContent(ptc)
	t.setupTableLayout(stop)
	t.setupEvents()
	return t
}

func (t *TablePanel) setupTableLayout(stop StopFunc) {
	t.Select(0, 0).
		SetFixed(1, 1).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				stop()
			}
			if key == tcell.KeyEnter {
				t.SetSelectable(true, true)
			}
		}).SetSelectedFunc(func(row int, column int) {
		t.GetCell(row, column).SetTextColor(tcell.ColorRed)
		t.SetSelectable(true, true)
	})
}

func (t *TablePanel) setupEvents() {
	t.SetSelectedFunc(func(row, col int) {
		t.onCellSelected.Notify(CellParams{row, col})
	})
}

type CellParams struct {
	Row, Col int
}

func (t *TablePanel) OnCellSelectedSubscribe(fn func(CellParams)) {
	t.onCellSelected.Subscribe(fn)
}
func (t *TablePanel) OnCellSelectedNotify(c CellParams) {
	t.onCellSelected.Notify(c)
}
