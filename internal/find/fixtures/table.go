// Package fixtures have test helpers.
package fixtures

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sahilm/fuzzy"
)

type Table struct {
	Keys   []string
	Values [][]string
}

func (t *Table) GetColumnCount() int {
	return len(t.Keys)
}

func (t *Table) GetColumn(col int) string {
	return t.Keys[col]
}
func (t *Table) GetRowCount() int {
	return len(t.Values)
}

func (t *Table) GetCell(row, col int) string {
	return t.Values[row][col]
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

// NewTable returns a fake kvSource.
// Each element is the index "r,c" as a string.
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

// FindSubstrings fake the fuzzy finder interface.
// Return the simple substring matches.
func FindSubstrings(pattern string, data []string) fuzzy.Matches {
	matches := fuzzy.Matches{}
	for i, s := range data {
		idx := strings.Index(s, pattern)
		if idx != -1 {
			matches = append(matches, newMatch(s, i, idx))
		}
	}
	return matches
}

func newMatch(word string, wordIdx int, chIdx int) fuzzy.Match {
	return fuzzy.Match{
		Str:            word,
		Index:          wordIdx,
		MatchedIndexes: []int{chIdx},
		Score:          1, // Score used to rank matches
	}
}
