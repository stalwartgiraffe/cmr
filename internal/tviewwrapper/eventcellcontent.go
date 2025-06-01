package tviewwrapper

import (
	"fmt"
	"time"

	"github.com/aarondl/opt/omitnull"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

type EventModelTextFunc func(e *gitlab.EventModel, projects map[int]gitlab.ProjectModel) string

type EventModelCellContent struct {
	title    string
	maxWidth int
	text     EventModelTextFunc
}

var defaultFormat = time.Stamp

func NewEventCellContents() []EventModelCellContent {
	return []EventModelCellContent{
		{
			title:    "ID",
			maxWidth: 8,
			text:     func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return fmt.Sprint(e.ID) },
		},
		{
			title:    "ProjectID",
			maxWidth: 20,
			text: func(e *gitlab.EventModel, projects map[int]gitlab.ProjectModel) string {
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
			title:    "AuthorUsername",
			maxWidth: 10,
			text:     func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return fmt.Sprint(e.AuthorUsername) },
		},
		{
			title:    "Title",
			maxWidth: 20,
			text:     func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.Title.GetOr("") },
		},
		{
			title:    "ActionName",
			maxWidth: 20,
			text:     func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return fmt.Sprint(e.ActionName) },
		},
		{
			title:    "TargetType",
			maxWidth: 20,
			text:     func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return fmt.Sprint(e.TargetType) },
		},
		{
			title:    "TargetTitle",
			maxWidth: 20,
			text:     func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.TargetTitle.GetOr("") },
		},
		{
			title:    "CreatedAt",
			maxWidth: 20,
			text: func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string {
				return fmt.Sprint(e.CreatedAt.Format(defaultFormat))
			},
		},
		{
			title:    "Data",
			maxWidth: 20,
			text:     func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.Data.GetOr("") },
		},
		// {
		// 	title:    "Author",
		// 	maxWidth: 30,
		// 	text: func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string {
		// 		p := e.Author
		// 		if p != nil {
		// 			return fmt.Sprintf("{%s}", p.String())
		// 		} else {
		// 			return ""
		// 		}
		// 	},
		// },

		{
			title:    "Imported",
			maxWidth: 10,
			text:     func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return getOrEmpty(e.Imported) },
		},
		{
			title:    "ImportedFrom",
			maxWidth: 10,
			text:     func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return getOrEmpty(e.ImportedFrom) },
		},
		{
			title:    "PushData",
			maxWidth: 10,
			text: func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string {
				if p, ok := e.PushData.Get(); ok {
					return p.String()
				}
				return ""
			},
		},

		{
			title:    "Note",
			maxWidth: 20,
			text:     func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.Note.String() },
		},
	}
}

func getOrEmpty[T any](b omitnull.Val[T]) string {
	if v, ok := b.Get(); ok {
		return fmt.Sprint(v)
	}
	return ""
}
