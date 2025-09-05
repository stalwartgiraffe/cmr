package merges

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	tw "github.com/stalwartgiraffe/cmr/internal/tviewwrapper"
)

type MergesRenderer interface {
	Run() error
}

var _ MergesRenderer = (*TuiMergesRenderer)(nil)

type TuiMergesRenderer struct {
	tviewApp *tview.Application
	stop     StopFn

	tablePage *tview.Flex

	filterPanel  *tw.BasicFilterPanel
	tablePanel   *tw.TablePanel
	detailsPanel *tw.TextDetailsPanel

	focusRing FocusRing
}

type FocusRing interface {
	Cycle(direction tw.RingDirection)
}

type StopFn func()

type MergesRepository interface {
	TableContents

	Load() error

	Filter(string)
}

type TableContents interface {
	GetRowCount() int
	GetColumnCount() int
	GetCell(row, col int) string
}

func NewTuiMergesRenderer(repo MergesRepository) *TuiMergesRenderer {
	tviewApp := tview.NewApplication()
	stop := tviewApp.Stop
	r := &TuiMergesRenderer{
		tviewApp:    tviewApp,
		tablePage:   tview.NewFlex(),
		filterPanel: tw.NewBasicFilterPanel(""),
		tablePanel: tw.NewTablePanel(
			tw.NewTwoBandTableContent(repo),
			stop),
		detailsPanel: tw.NewTextDetailsPanel(),
		stop:         stop,
	}

	r.focusRing = tw.NewFocusRing(tviewApp, r.filterPanel, r.tablePanel, r.detailsPanel)

	r.setupTablePage()
	r.setupKeyHandlers()
	r.setupEvents(repo)
	return r
}

func (r *TuiMergesRenderer) Run() error {
	return r.tviewApp.SetRoot(r.tablePage, true).SetFocus(r.tablePage).Run()
}

// setupTablePage laysout the panels in the page
func (r *TuiMergesRenderer) setupTablePage() {
	r.tablePage.SetDirection(tview.FlexRow)

	// row 1
	r.tablePage.AddItem(r.filterPanel, 3, 0, false)

	// row 2
	tableRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	tableRow.AddItem(r.tablePanel, 0, 2, true)                   // Table takes 2/3
	tableRow.AddItem(r.detailsPanel.GetPrimitive(), 0, 1, false) // Details takes 1/3
	r.tablePage.AddItem(tableRow, 0, 1, true)
}

// setupKeyHandlers configures keyboard navigation
func (r *TuiMergesRenderer) setupKeyHandlers() {
	r.tablePage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			r.stop()
			return nil
		case tcell.KeyTab:
			r.focusRing.Cycle(tw.NextDir)
			return nil
		case tcell.KeyBacktab:
			r.focusRing.Cycle(tw.PrevDir)
			return nil
		}
		return event
	})
}

func (r *TuiMergesRenderer) setupEvents(repo MergesRepository) {
	// Table selection handler
	r.tablePanel.SetSelectedFunc(func(row, col int) {
		if row > 0 { // Skip header row
			r.showTableRowDetails(row - 1)
		}
	})

	// Filter change handler
	r.filterPanel.OnChangeSubscribe(func(filterText string) {
		repo.Filter(filterText)
	})
}

// showTableRowDetails displays details for the selected table row
func (r *TuiMergesRenderer) showTableRowDetails(row int) {
	// This would need to be implemented by the specific table type
	// For now, just clear details
	//s := slices.Collect(maps.Values(requests))
	//details.ShowDetails(s[0])

	r.detailsPanel.Clear()
}

/*
// applyFilter applies the filter to the table content
func (r *TuiMergesRenderer) applyFilter(filterText string) {
	// This would need to be implemented to filter the table content
	// Implementation depends on the specific data being displayed
}
*/
