package find

import (
	"sort"
	"testing"

	mocks "github.com/stalwartgiraffe/cmr/internal/find/fixtures"

	"github.com/stretchr/testify/require"
)

func TestRemoveExcludedWords(t *testing.T) {
	setupSrc := func() TextTable {
		return &mocks.Table{
			Keys: []string{"Fruit"},
			Values: [][]string{
				{"Apple"},
				{"Banana"},
				{"Orange"},
				{"Grape"},
				{"Strawberry"},
				{"Mango"},
				{"Pineapple"},
			},
		}
	}

	tests := []struct {
		name     string
		kvSrc    TextTable
		excluded []int
		pattern  string
		want     []int
	}{
		{
			name:     "empty excluded",
			kvSrc:    setupSrc(),
			excluded: []int{},
			pattern:  "anana",
			want:     []int{},
		},
		{
			name:     "full excluded",
			kvSrc:    setupSrc(),
			excluded: []int{0, 1, 2, 3, 4, 5, 6},
			pattern:  "anana",
			want:     []int{0, 2, 3, 4, 5, 6},
		},
		{
			name:     "several excluded",
			kvSrc:    setupSrc(),
			excluded: []int{0, 1, 2, 3, 4, 5, 6},
			pattern:  "an",
			want:     []int{0, 3, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newFindSrc(tt.kvSrc)
			have := s.removeExcluded(tt.excluded, tt.pattern, 0)
			sort.Ints(have)
			require.Equal(t, tt.want, have)
		})
	}
}

func TwestRemoveExcludedSentences(t *testing.T) {
	setupSrc := func() TextTable {
		return &mocks.Table{
			Keys: []string{"Fruit"},
			Values: [][]string{
				{"The code, a tangled, thorny vine,"},      // 0
				{"A logic error, not a single line."},      // 1
				{"A bug unseen, a creeping dread,"},        // 2
				{"What's happening in the code instead?"},  // 3
				{"You step on through, the values check,"}, // 4
				{"A fix is found, no more a wreck."},       // 5
				{"The program runs, a happy end. "},        // 6
			},
		}
	}

	tests := []struct {
		name     string
		kvSrc    TextTable
		excluded []int
		pattern  string
		want     []int
	}{
			{
				name:     "empty excluded",
				kvSrc:    setupSrc(),
				excluded: []int{},
				pattern:  "anana",
				want:     []int{},
			},
			{
				name:     "full excluded",
				kvSrc:    setupSrc(),
				excluded: []int{0, 1, 2, 3, 4, 5, 6},
				pattern:  "dread",
				want:     []int{0, 2, 3, 4, 5, 6},
			},
		{
			name:     "several excluded",
			kvSrc:    setupSrc(),
			excluded: []int{0, 1, 2, 3, 4, 5, 6},
			pattern:  "no a",
			want:     []int{0, 2, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newFindSrc(tt.kvSrc)
			have := s.removeExcluded(tt.excluded, tt.pattern, 0)
			sort.Ints(have)
			require.Equal(t, tt.want, have)
		})
	}
}
