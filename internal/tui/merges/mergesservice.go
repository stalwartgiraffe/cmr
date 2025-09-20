package merges

import (
	"fmt"

	"github.com/stalwartgiraffe/cmr/events"
	"github.com/stalwartgiraffe/cmr/internal/find"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

type InMemoryMergesRepository struct {
	//projects  map[int]gitlab.ProjectModel
	//mergesMap gitlab.MergeRequestMap

	//mergeRequests views.DataView[gitlab.MergeRequestModel]
	//contents []MergeRequestModelContent

	table   *RecordTable
	view    *find.TableView
	changes events.Event[EmptyT]
}

type EmptyT = struct{}

type rowRef struct {
	idx int
	m   *gitlab.MergeRequestModel
}

func NewInMemoryMergesRepository() *InMemoryMergesRepository {
	return &InMemoryMergesRepository{}
}

func (r *InMemoryMergesRepository) Load() error {
	contents := NewMergeRequestContents()
	projects, err := gitlab.ReadProjects()
	if err != nil {
		return err
	}

	filepath := "ignore/my_recent_merge_request.yaml"
	mergesMap, err := gitlab.NewMergeRequestMapFromYaml(filepath)
	if err != nil {
		return err
	}
	r.table = NewRecordTable(contents, mergesMap, projects)
	r.view = find.NewTableView(r.table)
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
	//search = strings.ToLower(search)
	// r.mergeRequests.FilterAll(func(m *gitlab.MergeRequestModel) bool {
	// 	// pragmatically waste allocs for now
	// 	return strings.Contains(
	// 		strings.ToLower(getUserName(m)), search) ||
	// 		strings.Contains(strings.ToLower(m.Title), search)
	// })

	r.view.UpdateFind(search)
	r.changed()
}

func (r *InMemoryMergesRepository) GetRowCount() int {
	//return r.table.GetRowCount()
	return r.view.GetRowCount()
}

func (r *InMemoryMergesRepository) GetColumnCount() int {
	//return r.table.GetColumnCount()
	return r.view.GetColumnCount()
}

func (r *InMemoryMergesRepository) GetColumn(col int) string {
	//return r.table.GetColumn(col)
	return r.view.GetColumn(col)
}

func (r *InMemoryMergesRepository) GetCell(row int, col int) string {
	//return r.table.GetCell(row, col)
	return r.view.GetCell(row, col)
	/*
		content := r.contents[col]
		if row == 0 {
			return content.title
		}
		data := r.mergeRequests.Get(row - 1)
		return content.cell(data, r.projects)
	*/
}

func (r *InMemoryMergesRepository) GetRowRecord(row int) any {
	return r.table.records[row]
	/*
		if row == 0 {
			return nil
		}
		return r.mergeRequests.Get(row - 1)
	*/
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
