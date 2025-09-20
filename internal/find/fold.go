package find

import (
	"unicode"
	"unicode/utf8"
)

func asciiContainsAtFold(txt string, sub string, start int) (bool, int) {
	numStrBytes := len(txt)
	numSubBytes := len(sub)
	if numSubBytes == 0 {
		return true, 0
	}
	if numStrBytes == 0 {
		return false, 0
	}

	if (numStrBytes - start) < numSubBytes {
		return false, 0
	}
	b := 0
	for ; b < numSubBytes; b++ {
		if !byteEqualsFold(txt[start+b], sub[b]) {
			return false, 0
		}
	}
	return true, start + b
}

func utfContainsAtFold(txt string, sub string, start int) (bool, int) {
	numStrBytes := len(txt)
	numSubBytes := len(sub)
	if numSubBytes == 0 {
		return true, 0
	}
	if numStrBytes == 0 {
		return false, 0
	}

	if (numStrBytes - start) < numSubBytes {
		return false, 0
	}
	b := 0
	for b < numSubBytes {
		subRune, subWidth := utf8.DecodeRuneInString(sub[b:])
		if subRune == utf8.RuneError {
			return false, 0
		}
		strRune, strWidth := utf8.DecodeRuneInString(txt[start+b:])
		if strRune == utf8.RuneError {
			return false, 0
		}
		if subWidth != strWidth {
			return false, 0
		}
		if !utfEqualsFold(subRune, strRune) {
			return false, 0
		}
		b += subWidth
	}
	return true, start + b
}

func utfEqualsFold(a, b rune) bool {
	if a <= unicode.MaxASCII {
		return asciiEqualsFold(a, b)
	}
	return unicodeFoldEquals(a, b)
}

func unicodeFoldEquals(a, b rune) bool {
	if a == b {
		return true
	}
	if a < b {
		a, b = b, a
	}
	// General case. SimpleFold(x) returns the next equivalent rune > x
	// or wraps around to smaller values.
	r := unicode.SimpleFold(b)
	for r != b && r < a {
		r = unicode.SimpleFold(r)
	}
	return r == a
}

func asciiEqualsFold(a, b rune) bool {
	return byteEqualsFold(byte(a), byte(b))
}

// byteEqualsFold returns true if ascii byte characters are case insensitive equal
func byteEqualsFold(lo, up byte) bool {
	if lo == up {
		return true // just equal
	}
	// ascii fold
	if lo < up { // lexicographic sort, since upper_case < lower_case
		lo, up = up, lo // have to swap
	}
	return 'A' <= up && up <= 'Z' && // is_upper_case(up) &&
		lo == up+('a'-'A') // lower == to_lower(up)
}
