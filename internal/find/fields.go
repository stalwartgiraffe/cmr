// Package find searches stuff
package find

import (
	"sort"

	"github.com/sahilm/fuzzy"
)

type Fields struct {
}

func NewFields() *Fields {
	fields := &Fields{}

	return fields
}

type KVSource interface {
	NumKeys() int
	Key(col int) string
	NumValues() int
	Value(row, col int) string
}

/*
	return fuzzy.FindFrom(pattern, fs)
o
type Source interface {
	// The string to be matched at position i.
	String(i int) string
	// The length of the source. Typically is the length of the slice of things that you want to match.
	Len() int
type Matches []Match
type Match struct {
	// The matched string.
	Str string
	// The index of the matched string in the supplied slice.
	Index int
	// The indexes of matched characters. Useful for highlighting matches.
	MatchedIndexes []int
	// Score used to rank matches
	Score int
}
*/

func Find(rawPattern string, kvSrc KVSource) []int {
	terms := newTerms(rawPattern)
	keySource := newKeySource(kvSrc)
	numAllRows := kvSrc.NumValues()
	colSrc := newColumnSourceAllRows(kvSrc)
	searchedCols := colIdxSet{}

	for keyIdx, key := range terms.keys {
		matches := fuzzy.FindFromNoSort(key, keySource)
		if len(matches) < 0 {
			continue
		}

		for _, match := range matches {
			col := match.Index
			searchedCols[col] = empty{}
			colSrc = newColumnSource(kvSrc, col, colSrc)

			pattern := terms.keyPatterns[keyIdx]
			colSrc.removeMatches(pattern)
		}
	}

	// name description author
	for col := range kvSrc.NumKeys() {
		if _, ok := searchedCols[col]; ok {
			continue
		}
		// karl joe bob
		for _, pattern := range terms.valuePatterns {
			colSrc = newColumnSource(kvSrc, col, colSrc)
			colSrc.removeMatches(pattern)
		}
	}

	matchRows := subtractFromAll(colSrc.rows, numAllRows)
	return matchRows
}

type empty struct{}
type colIdxSet = map[int]empty

type keySource struct {
	src KVSource
}

func newKeySource(src KVSource) *keySource {
	return &keySource{
		src: src,
	}
}
func (s *keySource) String(col int) string {
	return s.src.Key(col)
}
func (s *keySource) Len() int {
	return s.src.NumKeys()
}

type columnSource struct {
	kvSrc  KVSource
	rows   []int
	column int

	findNoSort FindNoSortFn
}

type FindFn func(pattern string, data []string) fuzzy.Matches
type FindNoSortFn func(pattern string, data []string) fuzzy.Matches

//type FindFromFn func(pattern string, data fuzzy.Source) fuzzy.Matches
//type FindFromNoSortFn func(pattern string, data fuzzy.Source) fuzzy.Matches

func (s *columnSource) deepCopy() *columnSource {
	dst := make([]int, len(s.rows))
	copy(dst, s.rows)
	return &columnSource{
		kvSrc:  s.kvSrc,
		column: s.column,
		rows:   dst,

		findNoSort: fuzzy.FindNoSort,
	}
}
func newColumnSource(kvSrc KVSource, col int, colSrc *columnSource) *columnSource {
	return &columnSource{
		kvSrc:  kvSrc,
		column: col,
		rows:   colSrc.rows,
	}
}

func newColumnSourceAllRows(kvSrc KVSource) *columnSource {
	numAllRows := kvSrc.NumValues()
	rows := make([]int, numAllRows)
	for i := range numAllRows {
		rows[i] = i
	}
	return &columnSource{
		kvSrc: kvSrc,
		rows:  rows,
	}
}

func (s *columnSource) removeMatches(pattern string) {
	data := make([]string, 1)
	rows := s.rows
	i := 0
	for i < len(rows) {
		data[0] = s.kvSrc.Value(rows[i], s.column)
		m := s.findNoSort(pattern, data)
		if 0 < m.Len() { // match
			tail := len(rows) - 1
			rows[i] = rows[tail] // erase match
			rows = rows[:tail]   // clip tail
		} else {
			i++
		}
	}
}

// subtractFromAll returns the set inverse of src.
// src is assumed to be a unsorted slice of values in
// the range int [0..numAllRows), with no duplicates.
// src may be re-ordered
func subtractFromAll(src []int, numAllRows int) []int {
	numSrc := len(src)
	numDst := numAllRows - numSrc
	dst := make([]int, numDst)
	var s, d, rv int
	sort.Ints(src)
	for s < numSrc && d < numDst {
		sv := src[s]
		if rv == sv {
			rv++
			s++
		} else if rv < sv {
			dst[d] = rv
			d++
			rv++
		} else { // rv > sv
			s++
		}
	}
	fillVals(d, rv, dst)
	return dst
}

// fillVals starts at index i, and writes val into vals
// in ascending order.
func fillVals(i, val int, vals []int) {
	n := len(vals)
	for ; i < n; i++ {
		vals[i] = val
		val++
	}
}
