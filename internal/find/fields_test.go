package find

import (
	"testing"

	"github.com/stretchr/testify/require"

	mocks "github.com/stalwartgiraffe/cmr/internal/find/fixtures"
)

func TestColumnSource_removeMatches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		kvSrc        KVSource
		findNoSort   FindNoSortFn
		column       int
		initialCol   int
		initialRows  []int
		pattern      string
		mockMatches  map[string]bool // value -> should match
		expectedRows []int
	}{
		{
			name:         "no matches - all rows preserved",
			kvSrc:        mocks.NewTable(3, 2),
			findNoSort:   mocks.FindSubstrings,
			initialCol:   0,
			initialRows:  []int{0, 1, 2},
			pattern:      "xyz",
			expectedRows: []int{0, 1, 2},
		},
		{
			name:         "all matches - no rows preserved",
			kvSrc:        mocks.NewTable(3, 2),
			findNoSort:   mocks.FindSubstrings,
			initialCol:   0,
			initialRows:  []int{0, 1, 2},
			pattern:      "0",
			expectedRows: []int{},
		},
		{
			name:         "partial matches/middle removed",
			kvSrc:        mocks.NewTable(3, 2),
			findNoSort:   mocks.FindSubstrings,
			initialCol:   0,
			initialRows:  []int{2, 0, 1},
			pattern:      "1",
			expectedRows: []int{0, 2},
		},
		{
			name:         "partial matches/first removed",
			kvSrc:        mocks.NewTable(3, 2),
			findNoSort:   mocks.FindSubstrings,
			initialCol:   0,
			initialRows:  []int{2, 0, 1},
			pattern:      "2",
			expectedRows: []int{0, 1},
		},
		{
			name:         "partial matches/last removed",
			kvSrc:        mocks.NewTable(3, 2),
			findNoSort:   mocks.FindSubstrings,
			initialCol:   0,
			initialRows:  []int{2, 0, 1},
			pattern:      "1",
			expectedRows: []int{0, 2},
		},
		{
			name:         "empty rows",
			kvSrc:        mocks.NewTable(3, 2),
			findNoSort:   mocks.FindSubstrings,
			initialCol:   0,
			initialRows:  []int{},
			pattern:      "1",
			expectedRows: []int{},
		},
		{
			name:         "single row/no match",
			kvSrc:        mocks.NewTable(3, 2),
			findNoSort:   mocks.FindSubstrings,
			initialCol:   0,
			initialRows:  []int{0},
			pattern:      "1",
			expectedRows: []int{0},
		},
		{
			name:         "single row/one match",
			kvSrc:        mocks.NewTable(3, 2),
			findNoSort:   mocks.FindSubstrings,
			initialCol:   0,
			initialRows:  []int{0},
			pattern:      "0",
			expectedRows: []int{},
		},
		{
			name:         "multiple matches in sequence",
			kvSrc:        mocks.NewTable(20, 2),
			findNoSort:   mocks.FindSubstrings,
			initialCol:   0,
			initialRows:  []int{7, 8, 9, 10, 12},
			pattern:      "1",
			expectedRows: []int{7, 8, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			excluded := removeExcluded(
				tt.initialRows,
				mocks.FindSubstrings,
				tt.pattern,
				tt.kvSrc,
				0)

			require.ElementsMatch(t, tt.expectedRows, excluded)
		})
	}
}

func TestColumnSource_SubtractFromAll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		numAllRows   int
		numKeys      int
		rows         []int
		expectedRows []int
	}{
		{
			name:         "empty rows",
			numAllRows:   5,
			numKeys:      3,
			rows:         []int{},
			expectedRows: []int{0, 1, 2, 3, 4},
		},
		{
			name:         "single row excluded 2",
			numAllRows:   5,
			numKeys:      3,
			rows:         []int{2},
			expectedRows: []int{0, 1, 3, 4},
		},
		{
			name:         "single row excluded 0",
			numAllRows:   5,
			numKeys:      3,
			rows:         []int{0},
			expectedRows: []int{1, 2, 3, 4},
		},
		{
			name:         "multiple rows excluded",
			numAllRows:   5,
			numKeys:      3,
			rows:         []int{3, 1},
			expectedRows: []int{0, 2, 4},
		},
		{
			name:         "all rows excluded",
			numAllRows:   3,
			numKeys:      3,
			rows:         []int{0, 1, 2},
			expectedRows: []int{},
		},
		{
			name:         "unsorted rows excluded",
			numAllRows:   7,
			numKeys:      3,
			rows:         []int{5, 1, 3},
			expectedRows: []int{0, 2, 4, 6},
		},
		{
			name:         "consecutive rows excluded",
			numAllRows:   5,
			numKeys:      3,
			rows:         []int{1, 2, 3},
			expectedRows: []int{0, 4},
		},
		{
			name:         "first and last excluded",
			numAllRows:   5,
			numKeys:      3,
			rows:         []int{0, 4},
			expectedRows: []int{1, 2, 3},
		},
		{
			name:         "N equals 1, empty rows",
			numAllRows:   1,
			numKeys:      3,
			rows:         []int{},
			expectedRows: []int{0},
		},
		{
			name:         "N equals 1, one row excluded",
			numAllRows:   1,
			numKeys:      3,
			rows:         []int{0},
			expectedRows: []int{},
		},
		{
			name:         "N equals 0",
			numAllRows:   0,
			numKeys:      3,
			rows:         []int{},
			expectedRows: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := subtractFromAll(tt.rows, tt.numAllRows)
			require.Equal(t, tt.expectedRows, matches)
		})
	}
}

func TestEveryElement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		n        int
		expected []int
	}{
		{
			name: "zero elements",
			n:    0,
		},
		{
			name: "single element",
			n:    1,
		},
		{
			name: "multiple elements",
			n:    5,
		},
		{
			name: "ten elements",
			n:    10,
		},
		{
			name: "two elements",
			n:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := everyElement(tt.n)
			require.Equal(t, tt.n, len(result))
			// Verify each element equals its index
			for i, val := range result {
				require.Equal(t, i, val, "element at index %d should equal %d", i, i)
			}
		})
	}
}

func TestFillVals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		i        int
		val      int
		vals     []int
		expected []int
	}{
		{
			name:     "fill from beginning",
			i:        0,
			val:      10,
			vals:     make([]int, 5),
			expected: []int{10, 11, 12, 13, 14},
		},
		{
			name:     "fill from middle",
			i:        2,
			val:      20,
			vals:     []int{1, 2, 0, 0, 0},
			expected: []int{1, 2, 20, 21, 22},
		},
		{
			name:     "fill from end",
			i:        4,
			val:      100,
			vals:     []int{1, 2, 3, 4, 0},
			expected: []int{1, 2, 3, 4, 100},
		},
		{
			name:     "no fill needed - i equals length",
			i:        3,
			val:      50,
			vals:     []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "empty slice",
			i:        0,
			val:      5,
			vals:     []int{},
			expected: []int{},
		},
		{
			name:     "single element slice",
			i:        0,
			val:      42,
			vals:     make([]int, 1),
			expected: []int{42},
		},
		{
			name:     "negative starting value",
			i:        1,
			val:      -5,
			vals:     []int{99, 0, 0, 0},
			expected: []int{99, -5, -4, -3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fillVals(tt.i, tt.val, tt.vals)
			require.Equal(t, tt.expected, tt.vals)
		})
	}
}
