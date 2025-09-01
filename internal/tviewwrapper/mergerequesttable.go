package tviewwrapper

import (
	"fmt"
	"maps"
	"slices"
	"sort"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

type MergeRequestTextTable struct {
	mergeRequests []gitlab.MergeRequestModel
	projects      map[int]gitlab.ProjectModel
	contents      []MergeRequestModelContent
}

var _ TextTable = (*MergeRequestTextTable)(nil)

func NewMergeRequestTextTable(
	mrs gitlab.MergeRequestMap,
	projects map[int]gitlab.ProjectModel) *MergeRequestTextTable {

	s := slices.Collect(maps.Values(mrs))
	sort.Slice(s, func(i, j int) bool {
		return s[i].ID > s[j].ID
	})
	return &MergeRequestTextTable{
		mergeRequests: s,
		projects:      projects,
		contents:      NewMergeRequestContents(),
	}
}

func (t *MergeRequestTextTable) GetRowCount() int {
	return len(t.mergeRequests)
}

func (t *MergeRequestTextTable) GetColumnCount() int {
	return len(t.contents)
}

func (t *MergeRequestTextTable) GetCell(row int, col int) string {
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
