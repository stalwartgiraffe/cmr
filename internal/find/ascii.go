package find

type Char = byte
type Ascii []Char

func (txt Ascii) StartsWithFold(sub Ascii) bool {
	numTxt := len(txt)
	numSub := len(sub)
	if numTxt == 0 ||
		numSub == 0 ||
		numTxt < numSub {
		return false
	}
	for b := range numSub {
		if !CharEqualsFold(txt[b], sub[b]) {
			return false
		}
	}
	return true
}

const asciiToLower = 'a' - 'A'

// CharEqualsFold returns true if ascii byte characters are case insensitive equal
func CharEqualsFold(lo, up Char) bool {
	if lo == up {
		return true // just equal
	}
	// ascii fold
	if lo < up { // lexicographic sort, since upper_case < lower_case
		lo, up = up, lo // have to swap
	}
	return 'A' <= up && up <= 'Z' && // is_upper_case(up) &&
		lo == up+asciiToLower // to_lower(up)
}

func LetterToLower(up Char) Char {
	return up + asciiToLower
}
func LetterToUpper(lo Char) Char {
	return lo - asciiToLower
}
