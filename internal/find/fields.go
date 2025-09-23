// Package find searches stuff
package find

import (
	"sort"
	"strings"

	"github.com/stalwartgiraffe/cmr/internal/utils"
)

type TextTable interface {
	GetColumnCount() int
	GetColumn(col int) string
	GetRowCount() int
	GetCell(row, col int) string
}

type keyCols = map[string]int

func getColumnKeysToLower(src TextTable) keyCols {
	colMap := keyCols{}
	for col := range src.GetColumnCount() {
		colMap[strings.ToLower(src.GetColumn(col))] = col
	}
	return colMap
}

func Find(rawPattern string, kvSrc TextTable) []int {
	src := newFindSrc(kvSrc)
	patterns := newTerms(rawPattern)
	skipColumns := utils.Set[int]{}
	excluded := everyElement(src.numRows())
	excluded = src.removeAllMatches(excluded, skipColumns, patterns)
	rows := subtractFromAll(excluded, src.numRows())
	sort.Ints(rows)
	return rows
}
func FindOld(rawPattern string, kvSrc TextTable) []int {
	src := newFindSrc(kvSrc)
	patterns := newTerms(rawPattern)
	excluded, skipColumns := src.removeKeys(patterns)
	excluded = src.removeValues(excluded, skipColumns, patterns)
	rows := subtractFromAll(excluded, src.numRows())
	sort.Ints(rows)
	return rows
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

// everyElement returns the slice of 0..n-1
func everyElement(n int) []int {
	return fillVals(0, 0, make([]int, n))
}

// fillVals starts at index i, and writes val into vals
// in ascending order.
func fillVals(i, val int, vals []int) []int {
	n := len(vals)
	for ; i < n; i++ {
		vals[i] = val
		val++
	}
	return vals
}
