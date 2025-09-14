package find

import (
	"strings"

	"github.com/sahilm/fuzzy"

	"github.com/stalwartgiraffe/cmr/internal/utils"
)

type findSrc struct {
	kvSrc      KVSource
	findNoSort FindNoSortFn
}

func newFindSrc(kvSrc KVSource) *findSrc {
	return &findSrc{
		kvSrc:      kvSrc,
		findNoSort: fuzzy.FindNoSort,
	}
}

type FindFn func(pattern string, data []string) fuzzy.Matches
type FindNoSortFn func(pattern string, data []string) fuzzy.Matches

func (s *findSrc) numRows() int {
	return s.kvSrc.NumValues()
}
func (s *findSrc) removeKeys(patterns *terms) ([]int, utils.Set[int]) {
	keyMap := allKeyColsLower(s.kvSrc)
	excluded := everyElement(s.numRows())
	skipColumns := utils.Set[int]{}
	for kpIdx, rawKey := range patterns.keys {
		if col, ok := keyMap[strings.ToLower(rawKey)]; ok {
			skipColumns.Add(col)
			excluded = s.removeExcluded(
				excluded,
				patterns.keyPatterns[kpIdx],
				col,
			)
		}
	}
	return excluded, skipColumns
}

// removeValues removes the patterns that are in patterns, skipping columns in the skip list
func (s *findSrc) removeValues(excluded []int, skipColumns utils.Set[int], patterns *terms) []int {
	for col := range s.kvSrc.NumKeys() {
		if skipColumns.Contains(col) {
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
func (s *findSrc) removeExcluded(
	excluded []int,
	pattern string,
	col int) []int {
	value := []string{""}
	i := 0
	for i < len(excluded) {
		value[0] = s.kvSrc.Value(excluded[i], col)
		m := s.findNoSort(pattern, value)
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
