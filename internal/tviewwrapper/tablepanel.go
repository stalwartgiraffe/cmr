package tviewwrapper

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//type StopFunc func()

type TablePanel struct {
	*tview.Table
}

func NewTablePanel(ptc tview.TableContent, stop StopFunc) *TablePanel {
	p := &TablePanel{
		Table: tview.NewTable(),
	}
	p.Table.SetContent(ptc)
	layoutTablePanel(p.Table, stop)
	return p
}

func layoutTablePanel(table *tview.Table, stop StopFunc) {
	table.Select(0, 0).
		SetFixed(1, 1).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				stop()
			}
			if key == tcell.KeyEnter {
				table.SetSelectable(true, true)
			}
		}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(true, true)
	})
}
