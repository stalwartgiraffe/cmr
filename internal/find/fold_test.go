package find

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUtfContainsAtFold(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		sub       string
		runeStart int
		want      bool
	}{
		{
			name:      "test_case_name",
			str:       "watermelon",
			sub:       "water",
			runeStart: 0,
			want:      true,
		},
		{
			name:      "test_case_name",
			str:       "watermelon",
			sub:       "water",
			runeStart: 1,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have := utfContainsAtFold(tt.str, tt.sub, tt.runeStart)
			require.Equal(t, tt.want, have)
		})
	}
}
func TestUtfFoldEquals(t *testing.T) {
	tests := []struct {
		name string
		tr   rune
		sr   rune
		want bool
	}{
		// Same runes
		{
			name: "identical_lowercase",
			tr:   'a',
			sr:   'a',
			want: true,
		},
		{
			name: "identical_uppercase",
			tr:   'A',
			sr:   'A',
			want: true,
		},
		{
			name: "identical_digit",
			tr:   '5',
			sr:   '5',
			want: true,
		},
		// ASCII case-insensitive
		{
			name: "ascii_lower_upper",
			tr:   'a',
			sr:   'A',
			want: true,
		},
		{
			name: "ascii_upper_lower",
			tr:   'Z',
			sr:   'z',
			want: true,
		},
		{
			name: "ascii_middle_case",
			tr:   'm',
			sr:   'M',
			want: true,
		},
		// ASCII non-matching
		{
			name: "different_ascii_letters",
			tr:   'a',
			sr:   'b',
			want: false,
		},
		{
			name: "letter_vs_digit",
			tr:   'a',
			sr:   '1',
			want: false,
		},
		{
			name: "digit_vs_symbol",
			tr:   '5',
			sr:   '@',
			want: false,
		},
		// Unicode case folding
		{
			name: "unicode_german_ss",
			tr:   'ß',
			sr:   'S',
			want: false, // ß doesn't fold to S directly
		},
		{
			name: "unicode_accented_e",
			tr:   'é',
			sr:   'É',
			want: true,
		},
		{
			name: "unicode_greek_alpha",
			tr:   'α',
			sr:   'Α',
			want: true,
		},
		{
			name: "unicode_cyrillic_a",
			tr:   'а', // Cyrillic small a
			sr:   'А', // Cyrillic capital A
			want: true,
		},
		// Unicode non-matching
		{
			name: "different_unicode_chars",
			tr:   'α',
			sr:   'β',
			want: false,
		},
		{
			name: "unicode_vs_ascii",
			tr:   'α',
			sr:   'a',
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			have := utfEqualsFold(tt.tr, tt.sr)
			require.Equal(t, tt.want, have)
			have = unicodeFoldEquals(tt.tr, tt.sr)
			require.Equal(t, tt.want, have)
		})
	}
}

func TestAsciiEqualsFold(t *testing.T) {
	tests := []struct {
		name string
		tr   rune
		sr   rune
		want bool
	}{
		// Same runes
		{
			name: "identical_lowercase",
			tr:   'a',
			sr:   'a',
			want: true,
		},
		{
			name: "identical_uppercase",
			tr:   'A',
			sr:   'A',
			want: true,
		},
		{
			name: "identical_digit",
			tr:   '5',
			sr:   '5',
			want: true,
		},
		// ASCII case-insensitive
		{
			name: "ascii_lower_upper",
			tr:   'a',
			sr:   'A',
			want: true,
		},
		{
			name: "ascii_upper_lower",
			tr:   'Z',
			sr:   'z',
			want: true,
		},
		{
			name: "ascii_middle_case",
			tr:   'm',
			sr:   'M',
			want: true,
		},
		// ASCII non-matching
		{
			name: "different_ascii_letters",
			tr:   'a',
			sr:   'b',
			want: false,
		},
		{
			name: "letter_vs_digit",
			tr:   'a',
			sr:   '1',
			want: false,
		},
		{
			name: "digit_vs_symbol",
			tr:   '5',
			sr:   '@',
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			have := asciiEqualsFold(tt.tr, tt.sr)
			require.Equal(t, tt.want, have)
		})
	}
}
