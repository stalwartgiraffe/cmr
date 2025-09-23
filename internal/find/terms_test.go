package find

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTermsMatchValues(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		txt     string
		want    bool
	}{
		{
			name:    "one sentence",
			pattern: "want match",
			txt:     "this is some text i Want to match",
			want:    true,
		},
		// ASCII character matching
		{
			name:    "exact match lowercase",
			pattern: "hello",
			txt:     "hello world",
			want:    true,
		},
		{
			name:    "exact match uppercase",
			pattern: "HELLO",
			txt:     "HELLO WORLD",
			want:    true,
		},
		{
			name:    "case insensitive match",
			pattern: "hello",
			txt:     "HELLO world",
			want:    true,
		},
		{
			name:    "case insensitive mixed",
			pattern: "HeLLo",
			txt:     "hello WORLD",
			want:    true,
		},
		// Case-insensitive folding edge cases
		{
			name:    "all uppercase pattern in lowercase text",
			pattern: "WORLD",
			txt:     "hello world",
			want:    true,
		},
		{
			name:    "all lowercase pattern in uppercase text",
			pattern: "world",
			txt:     "HELLO WORLD",
			want:    true,
		},
		{
			name:    "mixed case folding",
			pattern: "WoRlD",
			txt:     "hello WORLD",
			want:    true,
		},
		// Mixed character types (letters, numbers, symbols)
		{
			name:    "alphanumeric match",
			pattern: "test123",
			txt:     "this is TEST123 data",
			want:    true,
		},
		{
			name:    "symbols and letters",
			pattern: "test@email",
			txt:     "user TEST@EMAIL address",
			want:    true,
		},
		{
			name:    "numbers only",
			pattern: "12345",
			txt:     "code 12345 found",
			want:    true,
		},
		{
			name:    "special characters",
			pattern: "hello-world",
			txt:     "say HELLO-WORLD today",
			want:    true,
		},
		{
			name:    "punctuation mix",
			pattern: "test.file",
			txt:     "found TEST.FILE here",
			want:    true,
		},
		// Edge cases
		{
			name:    "empty pattern",
			pattern: "",
			txt:     "any text",
			want:    true,
		},
		{
			name:    "empty text",
			pattern: "something",
			txt:     "",
			want:    false,
		},
		{
			name:    "both empty",
			pattern: "",
			txt:     "",
			want:    true,
		},
		{
			name:    "pattern longer than text",
			pattern: "verylongpattern",
			txt:     "short",
			want:    false,
		},
		{
			name:    "single character match",
			pattern: "a",
			txt:     "A",
			want:    true,
		},
		{
			name:    "single character no match",
			pattern: "a",
			txt:     "b",
			want:    false,
		},
		// Multiple value patterns
		{
			name:    "two patterns both found",
			pattern: "hello world",
			txt:     "say HELLO to the WORLD",
			want:    true,
		},
		{
			name:    "two patterns one missing",
			pattern: "hello missing",
			txt:     "say HELLO to the world",
			want:    false,
		},
		{
			name:    "three patterns all found",
			pattern: "one two three",
			txt:     "ONE and TWO and THREE",
			want:    true,
		},
		{
			name:    "three patterns middle missing",
			pattern: "one two three",
			txt:     "ONE and THREE",
			want:    false,
		},
		{
			name:    "overlapping patterns",
			pattern: "test testing",
			txt:     "TESTING is good",
			want:    true,
		},
		// Substring matching behavior
		{
			name:    "pattern at beginning",
			pattern: "hello",
			txt:     "HELLO world",
			want:    true,
		},
		{
			name:    "pattern at end",
			pattern: "world",
			txt:     "hello WORLD",
			want:    true,
		},
		{
			name:    "pattern in middle",
			pattern: "big",
			txt:     "the BIG house",
			want:    true,
		},
		{
			name:    "partial match should fail",
			pattern: "hello",
			txt:     "hell world",
			want:    false,
		},
		{
			name:    "pattern appears multiple times",
			pattern: "test",
			txt:     "TEST this TEST again",
			want:    true,
		},
		{
			name:    "no match anywhere",
			pattern: "missing",
			txt:     "hello world",
			want:    false,
		},
		// Order independence for multiple patterns
		{
			name:    "patterns found in different order",
			pattern: "hello world",
			txt:     "say HELLO to the WORLD",
			want:    true,
		},
		{
			name:    "sequential pattern matching order matters",
			pattern: "world hello",
			txt:     "say HELLO to the WORLD",
			want:    false,
		},
		{
			name:    "patterns with key-value mixed (should ignore key-value)",
			pattern: "?key:value hello",
			txt:     "say HELLO there",
			want:    true,
		},
		// Complex real-world scenarios
		{
			name:    "file path matching",
			pattern: "src main.go",
			txt:     "file SRC/MAIN.GO found",
			want:    true,
		},
		{
			name:    "email-like pattern",
			pattern: "user domain.com",
			txt:     "contact USER@DOMAIN.COM for help",
			want:    true,
		},
		{
			name:    "version number",
			pattern: "v1.2.3",
			txt:     "release V1.2.3 is ready",
			want:    true,
		},
		// Whitespace and special cases
		{
			name:    "pattern with internal spaces should not match",
			pattern: "hello world",
			txt:     "hellowrd",
			want:    false,
		},
		{
			name:    "consecutive patterns",
			pattern: "ab cd",
			txt:     "ABCD",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			search := newTerms(tt.pattern)
			have := search.matchValues(tt.txt)
			require.Equal(t, tt.want, have)
		})
	}
}
func TestNewTerms(t *testing.T) {
	tests := []struct {
		name             string
		rawPattern       string
		expectedKeys     []string
		expectedKeyPat   []string
		expectedValuePat []string
	}{
		{
			name:             "empty pattern",
			rawPattern:       "",
			expectedKeys:     []string{},
			expectedKeyPat:   []string{},
			expectedValuePat: []string{},
		},
		{
			name:             "single value pattern",
			rawPattern:       "test",
			expectedKeys:     []string{},
			expectedKeyPat:   []string{},
			expectedValuePat: []string{"test"},
		},
		{
			name:             "single key-value pattern",
			rawPattern:       "?key:value",
			expectedKeys:     []string{"key"},
			expectedKeyPat:   []string{"value"},
			expectedValuePat: []string{},
		},
		{
			name:             "mixed patterns",
			rawPattern:       "?name:john ?age:25 search",
			expectedKeys:     []string{"name", "age"},
			expectedKeyPat:   []string{"john", "25"},
			expectedValuePat: []string{"search"},
		},
		{
			name:             "multiple values",
			rawPattern:       "foo bar baz",
			expectedKeys:     []string{},
			expectedKeyPat:   []string{},
			expectedValuePat: []string{"foo", "bar", "baz"},
		},
		{
			name:             "invalid key pattern (no colon)",
			rawPattern:       "?key",
			expectedKeys:     []string{},
			expectedKeyPat:   []string{},
			expectedValuePat: []string{},
		},
		{
			name:             "invalid key pattern (empty key)",
			rawPattern:       "?:value",
			expectedKeys:     []string{},
			expectedKeyPat:   []string{},
			expectedValuePat: []string{},
		},
		{
			name:             "invalid key pattern (empty value)",
			rawPattern:       "?key:",
			expectedKeys:     []string{},
			expectedKeyPat:   []string{},
			expectedValuePat: []string{},
		},
		{
			name:             "complex mixed pattern",
			rawPattern:       "?status:active ?type:user admin search ?role:editor",
			expectedKeys:     []string{"status", "type", "role"},
			expectedKeyPat:   []string{"active", "user", "editor"},
			expectedValuePat: []string{"admin", "search"},
		},
		{
			name:             "bad separator 1",
			rawPattern:       "?status: active",
			expectedKeys:     []string{},
			expectedKeyPat:   []string{},
			expectedValuePat: []string{"active"},
		},
		{
			name:             "bad separator 2",
			rawPattern:       "?status :active",
			expectedKeys:     []string{},
			expectedKeyPat:   []string{},
			expectedValuePat: []string{":active"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newTerms(tt.rawPattern)

			require.Equal(t, tt.expectedKeys, result.keys, "keys mismatch")
			require.Equal(t, tt.expectedKeyPat, result.keyPatterns, "key patterns mismatch")
			require.Equal(t, tt.expectedValuePat, result.valuePatterns, "value patterns mismatch")
		})
	}
}
