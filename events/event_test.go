package events_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stalwartgiraffe/cmr/events"
)

func TestEvent_Subscribe(t *testing.T) {
	tests := []struct {
		name      string
		observers int
	}{
		{"single observer", 1},
		{"multiple observers", 3},
		{"no observers", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.Event[string]{}

			for i := 0; i < tt.observers; i++ {
				event.Subscribe(func(string) {})
			}

			require.Equal(t, tt.observers, event.Len())
		})
	}
}

func TestEvent_Notify(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected []string
	}{
		{"single notification", "test", []string{"test"}},
		{"empty string", "", []string{""}},
		{"unicode data", "🔥", []string{"🔥"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.Event[string]{}
			var received []string

			event.Subscribe(func(data string) {
				received = append(received, data)
			})

			event.Notify(tt.data)

			require.Equal(t, tt.expected, received)
		})
	}
}

func TestEvent_MultipleObservers(t *testing.T) {
	tests := []struct {
		name          string
		observerCount int
		data          int
	}{
		{"two observers", 2, 42},
		{"five observers", 5, 100},
		{"no observers", 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.Event[int]{}
			notifications := 0

			for i := 0; i < tt.observerCount; i++ {
				event.Subscribe(func(data int) {
					require.Equal(t, tt.data, data)
					notifications++
				})
			}

			event.Notify(tt.data)

			require.Equal(t, tt.observerCount, notifications)
		})
	}
}

func TestEvent_DifferentTypes(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"string type",
			func(t *testing.T) {
				event := &events.Event[string]{}
				var received string
				event.Subscribe(func(data string) { received = data })
				event.Notify("hello")
				require.Equal(t, "hello", received)
			},
		},
		{
			"int type",
			func(t *testing.T) {
				event := &events.Event[int]{}
				var received int
				event.Subscribe(func(data int) { received = data })
				event.Notify(123)
				require.Equal(t, 123, received)
			},
		},
		{
			"struct type",
			func(t *testing.T) {
				type TestStruct struct{ Value int }
				event := &events.Event[TestStruct]{}
				var received TestStruct
				event.Subscribe(func(data TestStruct) { received = data })
				event.Notify(TestStruct{Value: 456})
				require.Equal(t, TestStruct{Value: 456}, received)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
