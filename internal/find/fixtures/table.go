// Package fixtures have test helpers.
package fixtures

import (
	"fmt"
	"strconv"
)

type Table struct {
	Keys   []string
	Values [][]string
}

func (t *Table) NumKeys() int {
	return len(t.Keys)
}

func (t *Table) Key(col int) string {
	return t.Keys[col]
}

func (t *Table) NumValues() int {
	return len(t.Values)
}

func (t *Table) Value(row, col int) string {
	return t.Values[row][col]
}

func NewTable(rows, cols int) *Table {
	keys := make([]string, cols)
	for i := range cols {
		keys[i] = "key" + strconv.Itoa(i)
	}
	vals := make([][]string, rows)
	for r := range rows {
		vals[r] = make([]string, cols)
		for c := range cols {
			vals[r][c] = fmt.Sprintf("%d,%d", r, c)
		}
	}
	return &Table{
		Keys:   keys,
		Values: vals,
	}
}
