package merges

import (
	"fmt"
	"maps"
	"slices"
	"sort"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

type MergesRepository interface {
	TableContents

	Load() error

	// TODO remove
	GetCollections() (
		map[int]gitlab.ProjectModel,
		gitlab.MergeRequestMap,
	)
}

type InMemoryMergesRepository struct {
	projects      map[int]gitlab.ProjectModel
	mergesMap     gitlab.MergeRequestMap
	mergeRequests []gitlab.MergeRequestModel
	contents      []MergeRequestModelContent
}

var _ MergesRepository = (*InMemoryMergesRepository)(nil)

func NewInMemoryMergesRepository() *InMemoryMergesRepository {
	return &InMemoryMergesRepository{
		contents: NewMergeRequestContents(),
	}
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

	s := slices.Collect(maps.Values(r.mergesMap))
	sort.Slice(s, func(i, j int) bool {
		return s[i].ID > s[j].ID
	})

	r.mergeRequests = s
	return nil
}

// TODO remove this and push all behavior to repe
func (r *InMemoryMergesRepository) GetCollections() (
	map[int]gitlab.ProjectModel,
	gitlab.MergeRequestMap,
) {
	return r.projects, r.mergesMap
}

type TableContents interface {
	GetRowCount() int
	GetColumnCount() int
	GetCell(row, col int) string
}

func (t *InMemoryMergesRepository) GetRowCount() int {
	return len(t.mergeRequests)
}

func (t *InMemoryMergesRepository) GetColumnCount() int {
	return len(t.contents)
}

func (t *InMemoryMergesRepository) GetCell(row int, col int) string {
	content := t.contents[col]
	if row == 0 {
		return content.title
	}
	mr := &t.mergeRequests[row-1]
	return content.cell(mr, t.projects)
}

type MergeRequestModelContentFunc func(mr *gitlab.MergeRequestModel, projects map[int]gitlab.ProjectModel) string

type MergeRequestModelContent struct {
	title string
	cell  MergeRequestModelContentFunc
}

func NewMergeRequestContents() []MergeRequestModelContent {
	return []MergeRequestModelContent{
		{
			title: "ID",
			cell:  func(m *gitlab.MergeRequestModel, _ map[int]gitlab.ProjectModel) string { return fmt.Sprint(m.ID) },
		},
		{
			title: "ProjectID",
			cell: func(e *gitlab.MergeRequestModel, projects map[int]gitlab.ProjectModel) string {
				if p, ok := projects[e.ProjectID]; ok {
					return p.Name
				}
				return fmt.Sprint(e.ProjectID)
			},
		},
		{
			title: "AuthorUsername",
			cell: func(m *gitlab.MergeRequestModel, _ map[int]gitlab.ProjectModel) string {
				name := ""
				if m.Author != nil {
					name = m.Author.Username
				}
				return name
			},
		},
		{
			title: "Title",
			cell: func(m *gitlab.MergeRequestModel, _ map[int]gitlab.ProjectModel) string {
				return m.Title
			},
		},
	}
}
