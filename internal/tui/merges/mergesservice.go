package merges

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/stalwartgiraffe/cmr/events"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/views"
)

type InMemoryMergesRepository struct {
	projects  map[int]gitlab.ProjectModel
	mergesMap gitlab.MergeRequestMap

	mergeRequests views.DataView[gitlab.MergeRequestModel]
	contents      []MergeRequestModelContent

	changes events.Event[EmptyT]
}

type EmptyT = struct{}

type rowRef struct {
	idx int
	m   *gitlab.MergeRequestModel
}

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

	r.mergeRequests = views.NewDataView(s)
	r.mergeRequests.FilterAll(func(m *gitlab.MergeRequestModel) bool {
		return true
	})

	r.changed()
	return nil
}

func (r *InMemoryMergesRepository) Filter(search string) {
	// there is no simple standard library case insensitive string.Contains()
	// could use standard lib regex with i
	//
	// if perf becomes an concern
	//
	// consider third parties
	// github.com/charlievieth/strcase
	// import "github.com/sahilm/fuzzy"
	//
	// consider fancy alternative
	// pre generated indexed data structures - ie all txt in lower case
	search = strings.ToLower(search)
	r.mergeRequests.FilterAll(func(m *gitlab.MergeRequestModel) bool {
		// pragmatically waste allocs for now
		return strings.Contains(
			strings.ToLower(getUserName(m)), search) ||
			strings.Contains(strings.ToLower(m.Title), search)
	})
	r.changed()
}

func (r *InMemoryMergesRepository) GetRowCount() int {
	return r.mergeRequests.Len()
}

func (r *InMemoryMergesRepository) GetColumnCount() int {
	return len(r.contents)
}

func (r *InMemoryMergesRepository) GetCell(row int, col int) string {
	content := r.contents[col]
	if row == 0 {
		return content.title
	}
	data := r.mergeRequests.Get(row - 1)
	return content.cell(data, r.projects)
}

func (r *InMemoryMergesRepository) GetRowRecord(row int) any {
	if row == 0 {
		return nil
	}
	return r.mergeRequests.Get(row - 1)
}

type EmptyFn = func(EmptyT)

func (r *InMemoryMergesRepository) OnChanged(callback EmptyFn) {
	r.changes.Subscribe(callback)
}

func (r *InMemoryMergesRepository) changed() {
	r.changes.Notify(EmptyT{})
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
				return getUserName(m)
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

func getUserName(m *gitlab.MergeRequestModel) string {
	if m.Author != nil {
		return m.Author.Username
	}
	return ""
}
