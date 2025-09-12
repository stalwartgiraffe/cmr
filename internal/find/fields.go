// Package find searches stuff
package find

type Fields struct {
}

type Records interface {
	Keys() []string
	Values() [][]string
	Weights() []float64
}

func NewFields() *Fields {

	fields := &Fields{}

	return fields
}

type KVSource interface {
	NumKeys() int
	Key(col int) string
	NumValues() int
	Value(row, col int) string
}
