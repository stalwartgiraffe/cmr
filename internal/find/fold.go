package find


import (
	"unicode"
	"unicode/utf8"
)

func containsFoldAt(str string, sub string, runeStart int) bool {
	numStrBytes := len(str)
	numSubBytes := len(sub)
	if numStrBytes == 0 || numSubBytes == 0 {
		return numStrBytes == 0 && numSubBytes == 0
	}

	if (numStrBytes - runeStart) < numSubBytes {
		return false
	}
	i := 0
	for i < numSubBytes {
		subRune, subWidth := utf8.DecodeRuneInString(sub[i:])
		if subRune == utf8.RuneError {
			return false
		}
		strRune, strWidth := utf8.DecodeRuneInString(str[i:])
		if strRune == utf8.RuneError {
			return false
		}
		if subWidth != strWidth {
			return false
		}
		if !foldEquals(subRune, strRune) {
			return false
		}
		i += subWidth
	}
	return true
}


func foldEquals(a, b rune) bool {
	if a < utf8.RuneSelf {
		return asciiFoldEquals(a, b)
	}
	return utfFoldEquals(a, b)
}

func asciiFoldEquals(a, b rune) bool {
	if a == b {
		return true
	}
	if a < b {
		a, b = b, a
	}
	return 'A' <= b && b <= 'Z' && a == b+('a'-'A')
}

func utfFoldEquals(a, b rune) bool {
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

