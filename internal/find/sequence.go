package find

import (
	"unicode"
	"unicode/utf8"
)

// Taken from strings.EqualFold
func equalFold(tr, sr rune) bool {
	if tr == sr {
		return true
	}
	if tr < sr {
		tr, sr = sr, tr
	}
	// Fast check for ASCII.
	if tr < utf8.RuneSelf {
		// ASCII, and sr is upper case.  tr must be lower case.
		if 'A' <= sr && sr <= 'Z' && tr == sr+'a'-'A' {
			return true
		}
		return false
	}

	// General case. SimpleFold(x) returns the next equivalent rune > x
	// or wraps around to smaller values.
	r := unicode.SimpleFold(sr)
	for r != sr && r < tr {
		r = unicode.SimpleFold(r)
	}
	return r == tr
}
