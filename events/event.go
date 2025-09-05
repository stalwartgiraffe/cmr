// Package events provide the Observable pattern
package events

// Callback is a generic one argument delegate.
type Callback[T any] func(T)

// Event is a simple collection of callbacks that can be subscribed and notified.
// 
// All calls are synchronous. 
// No Unsubscribe is provided.
// Subscribers must not error.
// Callback order should not be assumed by subscribers.
type Event[T any] struct {
	callbacks []Callback[T]
}

// Subscribe adds a callback to the callbacks.
func (o *Event[T]) Subscribe(fn Callback[T]) {
	o.callbacks = append(o.callbacks, fn)
}

// Len returns the count of callbacks.
func (o *Event[T]) Len() int {
	return len(o.callbacks)
}

// Notify each callback with data.
func (o *Event[T]) Notify(data T) {
	for _, fn := range o.callbacks {
		fn(data)
	}
}
