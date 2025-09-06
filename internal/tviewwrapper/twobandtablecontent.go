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

type TwoBandTableContent struct {
	tview.TableContentReadOnly
	rowColors []tcell.Color
	table     TextTable
}

var _ tview.TableContent = (*TwoBandTableContent)(nil)

func NewTwoBandTableContent(table TextTable) *TwoBandTableContent {
	return &TwoBandTableContent{
		table: table,
		rowColors: []tcell.Color{
			tcell.ColorDarkBlue,
			tcell.ColorDarkSlateGrey,
			tcell.ColorBlack,
		},
	}
}

// GetRowCount return the count of rows in the table
func (c *TwoBandTableContent) GetRowCount() int {
	return c.table.GetRowCount()
}

// GetColumnCount returns the count of columns in the table.
func (c *TwoBandTableContent) GetColumnCount() int {
	return c.table.GetColumnCount()
}

// GetCell returns the contents of a table cell.
func (c *TwoBandTableContent) GetCell(row, col int) *tview.TableCell {
	cell := &tview.TableCell{
		Align:           tview.AlignLeft,
		Color:           tview.Styles.PrimaryTextColor,
		Transparent:     false, // must for false for BackgroundColor to be drawn
		Text:            c.table.GetCell(row, col),
		BackgroundColor: c.bandBackground(row),
	}
	return cell
}
func (c *TwoBandTableContent) bandBackground(row int) tcell.Color {
	if row == 0 {
		return c.rowColors[0]
	} else if ((row / 2) % 2) == 0 {
		return c.rowColors[1]
	} else {
		return c.rowColors[2]
	}
}
