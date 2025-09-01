package tviewwrapper

import (
	"fmt"
	"time"

	"github.com/aarondl/opt/omitnull"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
)

type EventModelTextFunc func(e *gitlab.EventModel, projects map[int]gitlab.ProjectModel) string

type EventModelCellContent struct {
	title string
	cell  EventModelTextFunc
}

var defaultFormat = time.Stamp

func NewEventCellContents() []EventModelCellContent {
	return []EventModelCellContent{
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
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return fmt.Sprint(e.AuthorUsername) },
		},
		{
			title: "Title",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return e.Title.GetOr("") },
		},
		{
			title: "ActionName",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return fmt.Sprint(e.ActionName) },
		},
		{
			title: "TargetType",
			cell:  func(e *gitlab.EventModel, _ map[int]gitlab.ProjectModel) string { return fmt.Sprint(e.TargetType) },
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
