// Package find searches stuff
package find

import "strings"

type terms struct {
	keys           []string
	keyPatterns    []string
	valuesPatterns []string
}

const KV_PREFIX = "?"
const KV_SEPERATOR = ":"

// newTerms accepts rawPattern of the form
// ?key1:val1 ?key2:val2 val3...
// and returns the parse terms
// whitespace is treated as separators.
func newTerms(rawPattern string) terms {
	rawTerms := strings.Fields(rawPattern)
	rawN := len(rawTerms)
	keyPatterns := make([]string, 0, rawN)
	keys := make([]string, 0, rawN)
	valuesPatterns := make([]string, 0, rawN)
	for _, term := range rawTerms {
		if strings.HasPrefix(term, KV_PREFIX) {
			keys, keyPatterns = parsekv(term, keys, keyPatterns)
		} else {
			valuesPatterns = append(valuesPatterns, term)
		}
	}
	return terms{keys, keyPatterns, valuesPatterns}
}

// parsekv accepts term in the form ?key:val and appends or discards
func parsekv(term string, keys []string, keyPatterns []string) ([]string, []string) {
	term = term[1:]
	idx := strings.Index(term, KV_SEPERATOR)
	if idx < 0 {
		return keys, keyPatterns
	}
	before, after := term[:idx], term[idx+1:]
	if len(before) < 1 || len(after) < 1 {
		return keys, keyPatterns
	}
	return append(keys, before), append(keyPatterns, after)
}
