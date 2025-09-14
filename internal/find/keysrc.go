package find

import "strings"

type keyCols = map[string]int

type keySource struct {
	src KVSource
}

func newKeySource(src KVSource) *keySource {
	return &keySource{
		src: src,
	}
}

func (s *keySource) String(col int) string {
	return s.src.Key(col)
}

func (s *keySource) Len() int {
	return s.src.NumKeys()
}

func allKeyColsLower(src KVSource) keyCols {
	set := keyCols{}
	for col := range src.NumKeys() {
		set[strings.ToLower(src.Key(col))] = col
	}
	return set
}
