package tviewwrapper

import (
	"maps"
	"slices"
	"sort"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

type EventsTextTable struct {
	events       []gitlab.EventModel
	projects     map[int]gitlab.ProjectModel
	cellContents []EventModelCellContent
}

var _ TextTable = (*EventsTextTable)(nil)

func NewEventsTextTable(
	events gitlab.EventMap,
	projects map[int]gitlab.ProjectModel) *EventsTextTable {
	s := slices.Collect(maps.Values(events))
	sort.Slice(s, func(i, j int) bool {
		return s[i].ID > s[j].ID
	})
	return &EventsTextTable{
		events:       s,
		projects:     projects,
		cellContents: NewEventCellContents(),
	}
}

func (t *EventsTextTable) GetRowCount() int {
	return len(t.events)
}
func (t *EventsTextTable) GetColumnCount() int {
	return len(t.cellContents)
}
func (t *EventsTextTable) GetCell(row int, col int) string {
	content := t.cellContents[col]
	if row == 0 {
		return content.title
	}
	event := &t.events[row-1]
	return content.cell(event, t.projects)
}
