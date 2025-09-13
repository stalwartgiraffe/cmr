package find

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColumnSource_SubtractFromAll(t *testing.T) {
	tests := []struct {
		name         string
		N            int
		rows         []int
		expectedRows []int
	}{
		{
			name:         "empty rows",
			N:            5,
			rows:         []int{},
			expectedRows: []int{0, 1, 2, 3, 4},
		},
		/*
		{
			name:         "single row excluded 2",
			N:            5,
			rows:         []int{2},
			expectedRows: []int{0, 1, 3, 4},
		},
		{
			name:         "single row excluded 0",
			N:            5,
			rows:         []int{0},
			expectedRows: []int{ 1, 2, 3, 4},
		},
		{
			name:         "multiple rows excluded",
			N:            6,
			rows:         []int{1, 3, 5},
			expectedRows: []int{0, 2, 4},
		},
		{
			name:         "all rows excluded",
			N:            3,
			rows:         []int{0, 1, 2},
			expectedRows: []int{},
		},
		{
			name:         "unsorted rows excluded",
			N:            7,
			rows:         []int{5, 1, 3},
			expectedRows: []int{0, 2, 4, 6},
		},
		{
			name:         "consecutive rows excluded",
			N:            5,
			rows:         []int{1, 2, 3},
			expectedRows: []int{0, 4},
		},
		{
			name:         "first and last excluded",
			N:            5,
			rows:         []int{0, 4},
			expectedRows: []int{1, 2, 3},
		},
		{
			name:         "N equals 1, empty rows",
			N:            1,
			rows:         []int{},
			expectedRows: []int{0},
		},
		{
			name:         "N equals 1, one row excluded",
			N:            1,
			rows:         []int{0},
			expectedRows: []int{},
		},
		{
			name:         "N equals 0",
			N:            0,
			rows:         []int{},
			expectedRows: []int{},
		},
		*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := &columnSource{
				N:    tt.N,
				rows: tt.rows,
			}
			
			col.subtractFromAll()
			
			require.Equal(t, tt.expectedRows, col.rows, "rows should match expected result")
		})
	}
}
