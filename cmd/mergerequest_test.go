package cmd

import (
	"embed"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stalwartgiraffe/cmr/internal/app"
)

func TestLastNIndex(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		count    int
		txt      string
		expected int
	}{
		{
			name: "single occurrence found",
			//---------01234567890123456789012345678
			body:     "hello world hello",
			count:    1,
			txt:      "hello",
			expected: 12,
		},
		{
			name: "multiple occurrences found",
			//---------01234567890123456789012345678
			body:     "hello world hello universe hello",
			count:    2,
			txt:      "hello",
			expected: 12,
		},
		{
			name: "all occurrences found",
			//---------01234567890123456789012345678
			body:     "hello world hello universe hello",
			count:    3,
			txt:      "hello",
			expected: 0,
		},
		{
			name:     "not enough occurrences",
			body:     "hello world hello",
			count:    3,
			txt:      "hello",
			expected: -1,
		},
		{
			name:     "text not found",
			body:     "hello world",
			count:    1,
			txt:      "foo",
			expected: -1,
		},
		{
			name:     "empty body",
			body:     "",
			count:    1,
			txt:      "hello",
			expected: -1,
		},
		{
			name:     "zero count",
			body:     "hello world hello",
			count:    0,
			txt:      "hello",
			expected: 17,
		},
		{
			name:     "single character search",
			body:     "a.b.c.d",
			count:    2,
			txt:      ".",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lastNIndex(tt.body, tt.count, tt.txt)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestNextNIndex(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		count    int
		txt      string
		expected int
	}{
		{
			name:     "single occurrence found",
			body:     "hello world hello",
			count:    1,
			txt:      "hello",
			expected: 1,
		},
		{
			name: "multiple occurrences found - second",
			//---------01234567890123
			body:     "hello world hello universe hello",
			count:    2,
			txt:      "hello",
			expected: 13,
		},
		{
			name: "multiple occurrences found - third",
			//---------01234567890123456789012345678
			body:     "hello world hello universe hello",
			count:    3,
			txt:      "hello",
			expected: 28,
		},
		{
			name:     "not enough occurrences",
			body:     "hello world hello",
			count:    3,
			txt:      "hello",
			expected: -1,
		},
		{
			name:     "text not found",
			body:     "hello world",
			count:    1,
			txt:      "foo",
			expected: -1,
		},
		{
			name:     "empty body",
			body:     "",
			count:    1,
			txt:      "hello",
			expected: -1,
		},
		{
			name:     "zero count",
			body:     "hello world hello",
			count:    0,
			txt:      "hello",
			expected: -1,
		},
		{
			name: "single character search",
			//---------01234567890123456789012345678
			body:     "a.b.c.d",
			count:    2,
			txt:      ".",
			expected: 4,
		},
		{
			name:     "overlapping pattern",
			body:     "aaaa",
			count:    2,
			txt:      "aa",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nextNIndex(tt.body, tt.count, tt.txt)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFindErrPosition(t *testing.T) {
	tests := []struct {
		name     string
		errtxt   string
		expected int
		isErr    bool
	}{
		{
			name:     "valid position found",
			errtxt:   "Error: The character at position 42 is invalid",
			expected: 42,
		},
		{
			name:     "position at start",
			errtxt:   "The character at position 0 causes error",
			expected: 0,
		},
		{
			name:     "position with extra whitespace",
			errtxt:   "Something wrong. The character at position    123   is bad",
			expected: 123,
		},
		{
			name:     "prefix not found",
			errtxt:   "Generic error message without position",
			expected: 0,
			isErr:    true,
		},
		{
			name:     "prefix found but no number",
			errtxt:   "The character at position abc is invalid",
			expected: 0,
			isErr:    true,
		},
		{
			name:     "prefix found but no text after",
			errtxt:   "The character at position",
			expected: 0,
			isErr:    true,
		},
		{
			name:     "empty string",
			errtxt:   "",
			expected: 0,
			isErr:    true,
		},
		{
			name:     "tread negative number as positive",
			errtxt:   "The character at position -5 is invalid",
			expected: 5,
		},
		{
			name:     "multiple numbers after prefix, get first one",
			errtxt:   "The character at position 10 20 is invalid",
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseNextInt(tt.errtxt)
			require.Equal(t, tt.expected, result)
			require.Equal(t, err != nil, tt.isErr)
		})
	}
}

func TestUnmarshalModels(t *testing.T) {
	appErr := app.NewAppErr()
	jsonBlob := readAllFile(t, "data/merge_requests.json")
	models, err := unmarshalModels(appErr.App, jsonBlob)
	require.NoError(t, err)
	require.NotNil(t, models)
}


//go:embed data/merge_requests.json
var loadTestsFS embed.FS

// readAllFile reads the contents of the embedded file
// if the embedded file is missing, its a compile error
func readAllFile(t *testing.T, name string) []byte {
	file, err := loadTestsFS.Open(name)
	require.NoError(t, err)
	defer file.Close()
	blob, err := io.ReadAll(file)
	require.NoError(t, err)
	return blob
}
