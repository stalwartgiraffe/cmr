package find

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	mocks "github.com/stalwartgiraffe/cmr/internal/find/fixtures"
)

func TestAllKeyCols(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		kvSrc       TextTable
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
			kvSrc: func() TextTable {
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

			result := getColumnKeysToLower(tt.kvSrc)

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

// TestAllKeyCols_EdgeCases tests edge cases and invariants
func TestAllKeyCols_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("result length equals number of keys", func(t *testing.T) {
		kvSrc := mocks.NewTable(5, 7)
		result := getColumnKeysToLower(kvSrc)

		require.Len(t, result, kvSrc.NumKeys())
	})

	t.Run("all columns are present", func(t *testing.T) {
		kvSrc := mocks.NewTable(3, 6)
		result := getColumnKeysToLower(kvSrc)

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
		result := getColumnKeysToLower(kvSrc)

		// Verify that result[kvSrc.Key(i)] == i for all valid i
		for i := 0; i < kvSrc.NumKeys(); i++ {
			key := kvSrc.Key(i)
			mappedCol, exists := result[key]
			require.True(t, exists, "key %s should exist in result", key)
			require.Equal(t, i, mappedCol, "key %s should map back to original column %d", key, i)
		}
	})
}
