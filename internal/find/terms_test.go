package find

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
