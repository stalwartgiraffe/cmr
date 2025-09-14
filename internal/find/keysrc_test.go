package find

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	mocks "github.com/stalwartgiraffe/cmr/internal/find/fixtures"
)

func TestNewKeySource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		kvSrc KVSource
	}{
		{
			name:  "empty table",
			kvSrc: mocks.NewTable(0, 0),
		},
		{
			name:  "single row, single column",
			kvSrc: mocks.NewTable(1, 1),
		},
		{
			name:  "multiple rows and columns",
			kvSrc: mocks.NewTable(5, 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := newKeySource(tt.kvSrc)

			require.NotNil(t, result)
			require.Equal(t, tt.kvSrc, result.src)
		})
	}
}

func TestKeySource_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		kvSrc       KVSource
		col         int
		expectedKey string
	}{
		{
			name:        "first column",
			kvSrc:       mocks.NewTable(3, 2),
			col:         0,
			expectedKey: "key0",
		},
		{
			name:        "second column",
			kvSrc:       mocks.NewTable(3, 2),
			col:         1,
			expectedKey: "key1",
		},
		{
			name:        "single column table",
			kvSrc:       mocks.NewTable(5, 1),
			col:         0,
			expectedKey: "key0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			keySource := newKeySource(tt.kvSrc)
			result := keySource.String(tt.col)

			require.Equal(t, tt.expectedKey, result)
		})
	}
}

func TestKeySource_Len(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		kvSrc       KVSource
		expectedLen int
	}{
		{
			name:        "empty table",
			kvSrc:       mocks.NewTable(0, 0),
			expectedLen: 0,
		},
		{
			name:        "single column",
			kvSrc:       mocks.NewTable(5, 1),
			expectedLen: 1,
		},
		{
			name:        "multiple columns",
			kvSrc:       mocks.NewTable(3, 4),
			expectedLen: 4,
		},
		{
			name:        "ten columns",
			kvSrc:       mocks.NewTable(2, 10),
			expectedLen: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			keySource := newKeySource(tt.kvSrc)
			result := keySource.Len()

			require.Equal(t, tt.expectedLen, result)
		})
	}
}

func TestAllKeyCols(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		kvSrc       KVSource
		expectedMap keyCols
	}{
		{
			name:        "empty table",
			kvSrc:       mocks.NewTable(0, 0),
			expectedMap: keyCols{},
		},
		{
			name:  "single column",
			kvSrc: mocks.NewTable(3, 1),
			expectedMap: keyCols{
				"key0": 0,
			},
		},
		{
			name:  "two columns",
			kvSrc: mocks.NewTable(3, 2),
			expectedMap: keyCols{
				"key0": 0,
				"key1": 1,
			},
		},
		{
			name:  "five columns",
			kvSrc: mocks.NewTable(2, 5),
			expectedMap: keyCols{
				"key0": 0,
				"key1": 1,
				"key2": 2,
				"key3": 3,
				"key4": 4,
			},
		},
		{
			name: "column names in caps",
			kvSrc: func() KVSource {
				src := mocks.NewTable(2, 5)
				for i := range src.Keys {
					src.Keys[i] = strings.ToUpper(src.Key(i))
				}
				return src
			}(),
			expectedMap: keyCols{
				"key0": 0,
				"key1": 1,
				"key2": 2,
				"key3": 3,
				"key4": 4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := allKeyColsLower(tt.kvSrc)

			require.Equal(t, tt.expectedMap, result)
			require.Len(t, result, len(tt.expectedMap))

			// Verify all mappings are correct
			for expectedKey, expectedCol := range tt.expectedMap {
				actualCol, exists := result[expectedKey]
				require.True(t, exists, "key %s should exist in result", expectedKey)
				require.Equal(t, expectedCol, actualCol, "key %s should map to column %d", expectedKey, expectedCol)
			}
		})
	}
}

// TestKeySource_Integration tests keySource as a fuzzy.Source interface
func TestKeySource_Integration(t *testing.T) {
	t.Parallel()

	// Create a test table with known keys
	kvSrc := mocks.NewTable(3, 4) // 3 rows, 4 columns (key0, key1, key2, key3)
	keySource := newKeySource(kvSrc)

	// Test interface compliance by using it like fuzzy.Source
	require.Equal(t, 4, keySource.Len())

	for i := 0; i < keySource.Len(); i++ {
		expectedKey := "key" + string(rune('0'+i))
		actualKey := keySource.String(i)
		require.Equal(t, expectedKey, actualKey)
	}
}

// TestAllKeyCols_EdgeCases tests edge cases and invariants
func TestAllKeyCols_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("result length equals number of keys", func(t *testing.T) {
		kvSrc := mocks.NewTable(5, 7)
		result := allKeyColsLower(kvSrc)

		require.Len(t, result, kvSrc.NumKeys())
	})

	t.Run("all columns are present", func(t *testing.T) {
		kvSrc := mocks.NewTable(3, 6)
		result := allKeyColsLower(kvSrc)

		// Check that all column indices from 0 to NumKeys-1 are present as values
		foundIndices := make(map[int]bool)
		for _, colIdx := range result {
			foundIndices[colIdx] = true
		}

		for i := 0; i < kvSrc.NumKeys(); i++ {
			require.True(t, foundIndices[i], "column index %d should be present", i)
		}
	})

	t.Run("bidirectional mapping consistency", func(t *testing.T) {
		kvSrc := mocks.NewTable(4, 5)
		result := allKeyColsLower(kvSrc)

		// Verify that result[kvSrc.Key(i)] == i for all valid i
		for i := 0; i < kvSrc.NumKeys(); i++ {
			key := kvSrc.Key(i)
			mappedCol, exists := result[key]
			require.True(t, exists, "key %s should exist in result", key)
			require.Equal(t, i, mappedCol, "key %s should map back to original column %d", key, i)
		}
	})
}

