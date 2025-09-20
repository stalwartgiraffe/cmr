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
			name:      "exact_match_at_start",
			str:       "watermelon",
			sub:       "water",
			runeStart: 0,
			want:      true,
		},
		{
			name:      "no_match_offset",
			str:       "watermelon",
			sub:       "water",
			runeStart: 1,
			want:      false,
		},
		{
			name:      "case_insensitive_match",
			str:       "WaterMelon",
			sub:       "water",
			runeStart: 0,
			want:      true,
		},
		{
			name:      "unicode_case_match",
			str:       "ΑΒΓΔΕ",
			sub:       "αβγ",
			runeStart: 0,
			want:      true,
		},
		// Empty string cases
		{
			name:      "both_empty",
			str:       "",
			sub:       "",
			runeStart: 0,
			want:      true,
		},
		{
			name:      "empty_str_nonempty_sub",
			str:       "",
			sub:       "test",
			runeStart: 0,
			want:      false,
		},
		{
			name:      "nonempty_str_empty_sub",
			str:       "test",
			sub:       "",
			runeStart: 0,
			want:      true,
		},
		// Length check failure
		{
			name:      "sub_longer_than_remaining",
			str:       "hello",
			sub:       "world",
			runeStart: 2,
			want:      false,
		},
		{
			name:      "sub_exactly_remaining_length",
			str:       "hello",
			sub:       "llo",
			runeStart: 2,
			want:      true,
		},
		// Character mismatch
		{
			name:      "character_mismatch",
			str:       "hello",
			sub:       "world",
			runeStart: 0,
			want:      false,
		},
		{
			name:      "partial_match_then_mismatch",
			str:       "helloworld",
			sub:       "help",
			runeStart: 0,
			want:      false,
		},
		// UTF-8 error cases - need invalid bytes to trigger RuneError
		{
			name:      "invalid_utf8_in_sub",
			str:       "hello",
			sub:       "\xff\xfe",
			runeStart: 0,
			want:      false,
		},
		{
			name:      "invalid_utf8_in_str",
			str:       "\xff\xfe",
			sub:       "hello",
			runeStart: 0,
			want:      false,
		},
		// Width mismatch case - this is tricky to trigger since same chars should have same width
		{
			name:      "mixed_width_non_matching",
			str:       "café", // é is 2 bytes
			sub:       "cafe", // e is 1 byte
			runeStart: 0,
			want:      false,
		},
		// Edge case: invalid UTF-8 in middle of substring
		{
			name:      "invalid_utf8_in_sub_middle",
			str:       "hello",
			sub:       "h\xff", // valid h followed by invalid byte
			runeStart: 0,
			want:      false,
		},
		// Specific case: valid sub but invalid str at matching position
		{
			name:      "valid_sub_invalid_str_at_position",
			str:       "h\xff", // h followed by invalid byte
			sub:       "he",    // valid substring
			runeStart: 0,
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
