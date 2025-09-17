package merges

import (
	"maps"
	"slices"
	"sort"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

type RecordTable struct {
	contents []MergeRequestModelContent
	records  []gitlab.MergeRequestModel
	projects map[int]gitlab.ProjectModel
}

func NewRecordTable(
	contents []MergeRequestModelContent,
	mergesMap gitlab.MergeRequestMap,
	projects map[int]gitlab.ProjectModel,
) *RecordTable {
	s := slices.Collect(maps.Values(mergesMap))
	sort.Slice(s, func(i, j int) bool {
		return s[i].ID > s[j].ID
	})

	return &RecordTable{
		contents: contents,
		records:  s,
		projects: projects,
	}
}
func (r *RecordTable) GetRowCount() int {
	return len(r.records)
}

func (r *RecordTable) GetColumnCount() int {
	return len(r.contents)
}

func (r *RecordTable) GetColumn(col int) string {
	return r.contents[col].title
}

func (r *RecordTable) GetCell(row int, col int) string {
	content := r.contents[col]
	record := &r.records[row]
	return content.cell(record, r.projects)
}
