package merges

import (
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

type MergesRepository interface {
	Load() error

	GetCollections() (
		map[int]gitlab.ProjectModel,
		gitlab.MergeRequestMap,
	)
}

type InMemoryMergesRepository struct {
	projects  map[int]gitlab.ProjectModel
	mergesMap gitlab.MergeRequestMap
}

func NewInMemoryMergesRepository() *InMemoryMergesRepository {
	return &InMemoryMergesRepository{}
}

func (r *InMemoryMergesRepository) Load() error {
	var err error
	r.projects, err = gitlab.ReadProjects()
	if err != nil {
		return err
	}

	filepath := "ignore/my_recent_merge_request.yaml"
	r.mergesMap, err = gitlab.NewMergeRequestMapFromYaml(filepath)
	if err != nil {
		return err
	}
	return nil
}

func (r *InMemoryMergesRepository) GetCollections() (
	map[int]gitlab.ProjectModel,
	gitlab.MergeRequestMap,
) {
	return r.projects, r.mergesMap
}
