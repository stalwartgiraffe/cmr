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
	patterns := newTerms(rawPattern)
	keySource := newKeySource(kvSrc)
	numAllRows := kvSrc.NumValues()
	skipColumns := idxSet{}
	colSrc := &columnSource{
		kvSrc: kvSrc,
	}

	excluded := everyElement(numAllRows)
	for keyIdx, key := range patterns.keys {
		matches := fuzzy.FindFromNoSort(key, keySource)
		if len(matches) < 0 {
			continue
		}

		for _, match := range matches {
			col := match.Index
			skipColumns[col] = empty{}
			colSrc = newColumnSource(kvSrc, col, colSrc)

			pattern := patterns.keyPatterns[keyIdx]
			excluded = colSrc.removeExcluded(
				excluded,
				pattern,
				col,
			)
		}
	}

	excluded = colSrc.removePatterns(excluded, skipColumns, patterns)

	return subtractFromAll(excluded, numAllRows)
}

type empty struct{}
type idxSet = map[int]empty

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

func (s *columnSource) removePatterns(excluded []int, skipColumns idxSet, patterns *terms) []int {
	for col := range s.kvSrc.NumKeys() {
		if _, ok := skipColumns[col]; ok {
			continue
		}
		for _, pattern := range patterns.valuePatterns {
			excluded = s.removeExcluded(
				excluded,
				pattern,
				col,
			)
		}
	}

	return excluded
}

// removeExcluded returns the rows which match pattern in src removed from excluded.
// The elements of the excluded slice may shuffled in place and the slice shortened.
func (s *columnSource) removeExcluded(
	excluded []int,
	pattern string,
	col int) []int {
	data := []string{""}
	i := 0
	for i < len(excluded) {
		data[0] = s.kvSrc.Value(excluded[i], col)

		m := s.findNoSort(pattern, data)
		if 0 < m.Len() { // on match, shuffle down last and pop
			last := len(excluded) - 1
			excluded[i] = excluded[last]
			excluded = excluded[:last]
		} else {
			i++
		}
	}
	return excluded
}

type FindFn func(pattern string, data []string) fuzzy.Matches
type FindNoSortFn func(pattern string, data []string) fuzzy.Matches

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

// everyElement returns the slice of 0..n-1
func everyElement(n int) []int {
	e := make([]int, n)
	for i := range n {
		e[i] = i
	}
	return e
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
