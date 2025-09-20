package find

import (
	"unicode"
	"unicode/utf8"
)

func utfContainsAtFold(str string, sub string, runeStart int) bool {
	numStrBytes := len(str)
	numSubBytes := len(sub)
	if numStrBytes == 0 || numSubBytes == 0 {
		return numStrBytes == 0 && numSubBytes == 0
	}

	if (numStrBytes - runeStart) < numSubBytes {
		return false
	}
	b := 0
	for b < numSubBytes {
		subRune, subWidth := utf8.DecodeRuneInString(sub[b:])
		if subRune == utf8.RuneError {
			return false
		}
		strRune, strWidth := utf8.DecodeRuneInString(str[runeStart+b:])
		if strRune == utf8.RuneError {
			return false
		}
		if subWidth != strWidth {
			return false
		}
		if !utfEqualsFold(subRune, strRune) {
			return false
		}
		b += subWidth
	}
	return true
}

func utfEqualsFold(a, b rune) bool {
	if a < utf8.RuneSelf {
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
	if a == b {
		return true
	}
	if a < b {
		a, b = b, a
	}
	return 'A' <= b && b <= 'Z' && a == b+('a'-'A')
}
