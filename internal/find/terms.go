// Package find searches stuff
package find

import "strings"

type terms struct {
	keys          []string
	keyPatterns   []string
	valuePatterns []string
}

const KVPrefix = "?"
const KVSeparator = ":"

// newTerms accepts rawPattern of the form
// ?key1:val1 ?key2:val2 val3...
// and returns the parse terms
// whitespace is treated as separators.
func newTerms(rawPattern string) *terms {
	rawTerms := strings.Fields(rawPattern)
	rawN := len(rawTerms)
	keyPatterns := make([]string, 0, rawN)
	keys := make([]string, 0, rawN)
	valuesPatterns := make([]string, 0, rawN)
	for _, term := range rawTerms {
		if strings.HasPrefix(term, KVPrefix) {
			keys, keyPatterns = parseKV(term, keys, keyPatterns)
		} else {
			valuesPatterns = append(valuesPatterns, term)
		}
	}
	return &terms{keys, keyPatterns, valuesPatterns}
}

// parseKV accepts term in the form ?key:val and appends
// parse errors are silently discarded
func parseKV(term string, keys []string, keyPatterns []string) ([]string, []string) {
	term = term[1:]
	idx := strings.Index(term, KVSeparator)
	if idx < 0 {
		return keys, keyPatterns
	}
	key, pattern := term[:idx], term[idx+1:]
	if len(key) < 1 || len(pattern) < 1 {
		return keys, keyPatterns
	}
	return append(keys, key), append(keyPatterns, pattern)
}

func (t *terms) matchValues(strTxt string) bool {
	txt := Ascii([]byte(strTxt))
	for _, v := range t.valuePatterns {
		val := Ascii([]byte(v))
		numValBytes := len(val)
		for !txt.StartsWithFold(val) {
			if len(txt) <= numValBytes {
				return false
			}
			// OPTIMIZE me with generic spans (start, len)
			txt = txt[1:]
		}
	}
	return true
}
