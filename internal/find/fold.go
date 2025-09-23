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
		if !CharEqualsFold(txt[start+b], sub[b]) {
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
	return CharEqualsFold(byte(a), byte(b))
}
