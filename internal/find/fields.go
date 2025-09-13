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

func Find(rawPattern string, kvSrc KVSource) {
	terms := newTerms(rawPattern)
	keySource := newKeySource(kvSrc)
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

	// name decrp author
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

	colSrc.subtractFromAll()

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
	N      int
}

func newColumnSource(kvSrc KVSource, col int, colSrc *columnSource) *columnSource {
	return &columnSource{
		kvSrc:  kvSrc,
		column: col,
		N:      colSrc.N,
		rows:   colSrc.rows,
	}
}

func newColumnSourceAllRows(kvSrc KVSource) *columnSource {
	N := kvSrc.NumValues()
	rows := make([]int, N)
	for i := range N {
		rows[i] = i
	}
	return &columnSource{
		kvSrc: kvSrc,
		N:     N,
		rows:  rows,
	}
}

func (s *columnSource) removeMatches(pattern string) {
	data := make([]string, 1)
	rows := s.rows
	i := 0
	for i < len(rows) {
		data[0] = s.kvSrc.Value(rows[i], s.column)
		m := fuzzy.FindNoSort(pattern, data)
		if 0 < m.Len() { // match
			tail := len(rows) - 1
			rows[i] = rows[tail] // remove match
			rows = rows[:tail]   // remove tail
		} else {
			i++
		}
	}
}

func (col *columnSource) subtractFromAll() {
	src := col.rows
	sort.Ints(src)
	numSrc := len(col.rows)
	numDiff := col.N - numSrc
	diff := make([]int, numDiff)
	s := 0
	d := 0
	rv := 0
	for s < numSrc && d < numDiff {
		sv := src[s]
		if rv == sv {
			rv++
			s++
		} else if rv < sv {
			diff[d] = rv
			d++
			rv++
		} else { // rv > sv
			s++
		}
	}
	for d < numDiff {
		diff[d] = rv
		d++
		rv++
	}
}
