package tviewwrapper

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"time"

	"github.com/aarondl/opt/omitnull"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

type EventsTextTable struct {
	events   []gitlab.EventModel
	projects map[int]gitlab.ProjectModel
	contents []EventModelContent
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
		events:   s,
		projects: projects,
		contents: NewEventContents(),
	}
}

func (t *EventsTextTable) GetRowCount() int {
	return len(t.events)
}
func (t *EventsTextTable) GetColumnCount() int {
	return len(t.contents)
}
func (t *EventsTextTable) GetCell(row int, col int) string {
	content := t.contents[col]
	if row == 0 {
		return content.title
	}
	event := &t.events[row-1]
	return content.cell(event, t.projects)
}

type EventModelTextFunc func(e *gitlab.EventModel, projects map[int]gitlab.ProjectModel) string

type EventModelContent struct {
	title string
	cell  EventModelTextFunc
}

var defaultFormat = time.Stamp

func NewEventContents() []EventModelContent {
	return []EventModelContent{
		{
			title: "ID",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return fmt.Sprint(e.ID) },
		},
		{
			title: "ProjectID",
			cell: func(e *gitlab.EventModel, projects map[int]gitlab.ProjectModel) string {
				if p, ok := projects[e.ProjectID]; ok {
					return p.Name
				}
				return fmt.Sprint(e.ProjectID)
			},
		},
		/*
				{
					title:    "TargetID",
					maxWidth: 5,
					text:     func(_ map[int]gitlab.ProjectModel, e *gitlab.EventModel) string { return fmt.Sprint(e.TargetID) },
				},
			{
				title:    "AuthorID",
				maxWidth: 5,
				text:     func(_ map[int]gitlab.ProjectModel, e *gitlab.EventModel) string { return fmt.Sprint(e.AuthorID) },
			},
		*/
		{
			title: "AuthorUsername",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.AuthorUsername },
		},
		{
			title: "Title",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.Title.GetOr("") },
		},
		{
			title: "ActionName",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.ActionName },
		},
		{
			title: "TargetType",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.TargetType },
		},
		{
			title: "TargetTitle",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.TargetTitle.GetOr("") },
		},
		{
			title: "CreatedAt",
			cell: func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string {
				return fmt.Sprint(e.CreatedAt.Format(defaultFormat))
			},
		},
		{
			title: "Data",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.Data.GetOr("") },
		},

		{
			title: "Imported",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return getOrEmpty(e.Imported) },
		},
		{
			title: "ImportedFrom",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return getOrEmpty(e.ImportedFrom) },
		},
		{
			title: "PushData",
			cell: func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string {
				if p, ok := e.PushData.Get(); ok {
					return p.String()
				}
				return ""
			},
		},
		{
			title: "Note",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.Note.String() },
		},
	}
}

func getOrEmpty[T any](b omitnull.Val[T]) string {
	if v, ok := b.Get(); ok {
		return fmt.Sprint(v)
	}
	return ""
}
