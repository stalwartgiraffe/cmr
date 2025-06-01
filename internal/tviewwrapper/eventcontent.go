package tviewwrapper

import (
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"

	"golang.org/x/exp/maps"
)

type EventsContent struct {
	tview.TableContentReadOnly

	events   []gitlab.EventModel
	projects map[int]gitlab.ProjectModel

	cellContents []EventModelCellContent

	rowColors []tcell.Color
}

var _ tview.TableContent = (*EventsContent)(nil)

func NewEventsContent(events gitlab.EventMap, projects map[int]gitlab.ProjectModel) *EventsContent {
	s := maps.Values(events)
	sort.Slice(s, func(i, j int) bool {
		return s[i].ID > s[j].ID
	})
	return &EventsContent{
		events:       s,
		projects:     projects,
		cellContents: NewEventCellContents(),
		rowColors: []tcell.Color{
			tcell.ColorDarkBlue,
			tcell.ColorDarkSlateGrey,
			tcell.ColorBlack,
		},
	}
}

func (c *EventsContent) GetCell(row, col int) *tview.TableCell {
	width, bgColor, txt := c.renderCell(row, col)
	cell := &tview.TableCell{
		Text:            txt,
		Align:           tview.AlignLeft,
		Color:           tview.Styles.PrimaryTextColor,
		Transparent:     false, // must for false for BackgroundColor to be drawn
		BackgroundColor: bgColor,
		MaxWidth:        width,
	}
	return cell
}

func (c *EventsContent) renderCell(row, col int) (
	maxWidth int,
	bgColor tcell.Color,
	txt string,
) {
	bgColor = c.rowColors[rowBackground(row)]
	content := c.cellContents[col]
	if row == 0 {
		txt = content.title
	} else {
		event := &c.events[row]
		txt = content.text(event, c.projects)
	}
	return
}

const (
	titleIdx = iota
	firstBandIdx
	secondBandIdx
)

func rowBackground(row int) int {
	if row == 0 {
		return titleIdx
	} else if ((row / 2) % 2) == 0 {
		return firstBandIdx
	} else {
		return secondBandIdx
	}
}

// Return the total number of rows in the table.
func (c *EventsContent) GetRowCount() int {
	return len(c.events)
}

// Return the total number of columns in the table.
func (c *EventsContent) GetColumnCount() int {
	return len(c.cellContents)
}
