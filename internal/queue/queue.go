package queue

// Keeping below as var so it is possible to run the slice size bench tests with no coding changes.
var (
	// firstSliceSize holds the size of the first slice.
	firstSliceSize = 1

	// maxFirstSliceSize holds the maximum size of the first slice.
	maxFirstSliceSize = 16

	// maxInternalSliceSize holds the maximum size of each internal slice.
	maxInternalSliceSize = 128
)

// Que represents an unbounded, dynamically growing FIFO queue.
// The zero value for queue is an empty queue ready to use.
type Que[T any] struct {
	// head points to the first node of the linked list.
	head *page[T]

	// tail points to the last node of the linked list.
	// In an empty queue, head and tail points to the same node.
	tail *page[T]

	// first is the index pointing to the current first element in the queue
	// (i.e. first element added in the current queue values).
	first int

	// Len holds the current queue values length.
	length int

	// lastSliceSize holds the size of the last created internal slice.
	lastSliceSize int
}

// page represents a queue page.
// Each page holds a slice of user managed values.
type page[T any] struct {
	// vals holds the list of user added values in this node.
	vals []T

	// nexth points to the next node in the linked list.
	next *page[T]
}

// newPage returns an initialized node.
func newPage[T any](capacity int) *page[T] {
	return &page[T]{
		vals: make([]T, 0, capacity),
	}
}

// New returns an initialized queue.
func New[T any]() *Que[T] {
	return &Que[T]{}
}

// Len returns the number of elements of queue q.
// The complexity is O(1).
func (q *Que[T]) Len() int {
	return q.length
}

// Front returns the first element of queue q or nil if the queue is empty.
// The second, bool result indicates whether a valid value was returned;
//
//	if the queue is empty, false will be returned.
//
// The complexity is O(1).
func (q *Que[T]) Front() (T, bool) {
	if q.head == nil {
		var zero T
		return zero, false
	}
	return q.head.vals[q.first], true
}

// Push adds a value to the queue.
// The complexity is O(1).
func (q *Que[T]) Push(v T) {
	if q.head == nil {
		h := newPage[T](firstSliceSize)
		q.head = h
		q.tail = h
		q.lastSliceSize = maxFirstSliceSize
	} else if len(q.tail.vals) >= q.lastSliceSize {
		n := newPage[T](maxInternalSliceSize)
		q.tail.next = n
		q.tail = n
		q.lastSliceSize = maxInternalSliceSize
	}

	q.tail.vals = append(q.tail.vals, v)
	q.length++
}

// Pop retrieves and removes the current element from the queue.
// The second, bool result indicates whether a valid value was returned;
//
//	if the queue is empty, false will be returned.
//
// The complexity is O(1).
func (q *Que[T]) Pop() (T, bool) {
	var zero T
	if q.head == nil {
		return zero, false
	}

	v := q.head.vals[q.first]
	q.head.vals[q.first] = zero // Avoid memory leaks
	q.length--
	q.first++
	if q.first >= len(q.head.vals) {
		n := q.head.next
		q.head.next = nil // Avoid memory leaks
		q.head = n
		q.first = 0
	}
	return v, true
}
