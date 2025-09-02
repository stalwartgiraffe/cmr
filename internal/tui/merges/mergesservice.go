package merges

import (
	"maps"
	"slices"
	"sort"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

type MergesRepository interface {
	Load() error
}

type InMemoryMergesRepository struct {
	merges []gitlab.MergeRequestModel
}

func (r *InMemoryMergesRepository) Load() error {
	filepath := "ignore/my_recent_merge_request.yaml"
	mergesMap, err := gitlab.NewMergeRequestMapFromYaml(filepath)
	if err != nil {
		return err
	}
	r.index(mergesMap)
	return nil
}

func (r *InMemoryMergesRepository) index(mergesMap gitlab.MergeRequestMap) {
	merges := slices.Collect(maps.Values(mergesMap))
	sort.Slice(merges, func(i, j int) bool {
		return merges[i].ID > merges[j].ID
	})
	r.merges = merges
}
