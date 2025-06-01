package gitlab

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestEventMap_ProjectIDs(t *testing.T) {
	g := NewWithT(t)

	cases := map[string]struct {
		events EventMap

		want []int
	}{
		"empty": {
			want: make([]int, 0),
		},
		"one value": {
			events: EventMap{
				1: {
					ProjectID: 13,
				},
			},
			want: []int{13},
		},
		"two value": {
			events: EventMap{
				1: {
					ProjectID: 13,
				},
				2: {
					ProjectID: 11,
				},
			},
			want: []int{11, 13},
		},
		"two value, one dup": {
			events: EventMap{
				1: {
					ProjectID: 13,
				},
				2: {
					ProjectID: 11,
				},
				3: {
					ProjectID: 11,
				},
			},
			want: []int{11, 13},
		},
		"two value, two dup": {
			events: EventMap{
				1: {
					ProjectID: 13,
				},
				2: {
					ProjectID: 11,
				},
				3: {
					ProjectID: 11,
				},
				4: {
					ProjectID: 13,
				},
			},
			want: []int{11, 13},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			g.Expect(tc.want).To(Equal(tc.events.ProjectIDs()))
		})
	}
}
