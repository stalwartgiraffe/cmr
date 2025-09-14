package utils

// Set is a generic set type.
type empty = struct{}
type Set[T comparable] map[T]empty

// NewSet creates a new set.
func NewSet[T comparable](items ...T) Set[T] {
	s := make(Set[T])
	s.Add(items...)
	return s
}

// Add adds elements to the set.
func (s Set[T]) Add(items ...T) {
	for _, item := range items {
		s[item] = empty{}
	}
}

// Contains checks if an element is in the set.
func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

// Remove removes an element from the set.
func (s Set[T]) Remove(item T) {
	delete(s, item)
}
