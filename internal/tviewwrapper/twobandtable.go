package tviewwrapper

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TextTable interface {
	GetRowCount() int
	GetColumnCount() int
	GetCell(row, col int) string
}

type TwoBandTable struct {
	tview.TableContentReadOnly
	rowColors []tcell.Color
	table     TextTable
}

var _ tview.TableContent = (*TwoBandTable)(nil)

func NewTwoBandTable(table TextTable) *TwoBandTable {
	return &TwoBandTable{
		table: table,
		rowColors: []tcell.Color{
			tcell.ColorDarkBlue,
			tcell.ColorDarkSlateGrey,
			tcell.ColorBlack,
		},
	}
}

// GetRowCount return the count of rows in the table
func (c *TwoBandTable) GetRowCount() int {
	return c.table.GetRowCount()
}

// GetColumnCount returns the count of columns in the table.
func (c *TwoBandTable) GetColumnCount() int {
	return c.table.GetColumnCount()
}

// GetCell returns the contents of a table cell.
func (c *TwoBandTable) GetCell(row, col int) *tview.TableCell {
	cell := &tview.TableCell{
		Align:           tview.AlignLeft,
		Color:           tview.Styles.PrimaryTextColor,
		Transparent:     false, // must for false for BackgroundColor to be drawn
		Text:            c.table.GetCell(row, col),
		BackgroundColor: c.bandBackground(row),
		//MaxWidth:        width,
	}
	return cell
}
func (c *TwoBandTable) bandBackground(row int) tcell.Color {
	if row == 0 {
		return c.rowColors[0]
	} else if ((row / 2) % 2) == 0 {
		return c.rowColors[1]
	} else {
		return c.rowColors[2]
	}
}
