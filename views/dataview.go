// Package views provides data views of records
package views

type DataView[T any] struct {
	records    []T
	filterView []Ref[T]
}

type Ref[T any] struct {
	Idx  int
	Data *T
}

type isMatchFn[T any] = func(*T) bool

func NewDataView[T any](records []T) DataView[T] {
	return DataView[T]{
		records:    records,
		filterView: make([]Ref[T], 0, len(records)),
	}
}

func (v *DataView[T]) FilterAll(isMatch isMatchFn[T]) []Ref[T] {
	v.filterView = matchAll(v.records, v.filterView, isMatch)
	return v.filterView
}

func (v *DataView[T]) Len() int {
	return len(v.filterView)
}

func (v *DataView[T]) Get(i int) *T {
	return v.filterView[i].Data
}

func matchAll[T any](records []T, view []Ref[T], isMatch isMatchFn[T]) []Ref[T] {
	view = view[:0]
	for i := range records {
		p := &records[i]
		if isMatch(p) {
			view = append(view, Ref[T]{i, p})
		}
	}
	return view
}
