package tviewwrapper

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StopFunc func()

func MakeContentTable(ptc tview.TableContent, stop StopFunc) *tview.Table {
	table := tview.NewTable()

	table.SetContent(ptc)
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
	return table
}
