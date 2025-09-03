package merges

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/stalwartgiraffe/cmr/internal/tviewwrapper"
)

type MergesRenderer interface {
	Run() error
}

var _ MergesRenderer = (*TuiMergesRenderer)(nil)

type TuiMergesRenderer struct {
	tviewApp *tview.Application
	stop     StopFn

	tablePage *tview.Flex

	filter     *tviewwrapper.BasicFilter
	tablePanel *tviewwrapper.TablePanel
	details    *tviewwrapper.TextDetails

	currentFocus int
	focusOps     []focusOp
	panels       []tview.Primitive
}

type StopFn func()

func NewTuiMergesRenderer(repo MergesRepository) *TuiMergesRenderer {
	tviewApp := tview.NewApplication()
	stop := tviewApp.Stop
	r := &TuiMergesRenderer{
		tviewApp:  tviewApp,
		tablePage: tview.NewFlex(),

		filter: tviewwrapper.NewBasicFilter(""),

		tablePanel: tviewwrapper.NewTablePanel(
			tviewwrapper.NewTwoBandTableContent(repo),
			stop),

		details: tviewwrapper.NewTextDetails(),
		stop:    stop,
	}

	r.panels = []tview.Primitive{
		r.filter,
		r.tablePanel,
		r.details,
	}

	r.setupLayout()
	r.setupFocusRing()
	r.setupKeyHandlers()
	return r
}

func (r *TuiMergesRenderer) Run() error {
	return r.tviewApp.SetRoot(r.tablePage, true).SetFocus(r.tablePage).Run()
}

func (r *TuiMergesRenderer) setupLayout() {
	// Main horizontal layout
	r.tablePage.SetDirection(tview.FlexRow)
	r.tablePage.AddItem(r.filter, 3, 0, false)

	tableRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	tableRow.AddItem(r.tablePanel, 0, 2, true)              // Table takes 2/3
	tableRow.AddItem(r.details.GetPrimitive(), 0, 1, false) // Details takes 1/3
	r.tablePage.AddItem(tableRow, 0, 1, true)
}

type focusOp struct {
	focus func()
}

func (r *TuiMergesRenderer) setupFocusRing() {
	r.focusOps = []focusOp{
		{
			focus: func() {
				r.tviewApp.SetFocus(r.tablePanel)
			},
		},
		{
			focus: func() {
				r.tviewApp.SetFocus(r.details.GetPrimitive())
			},
		},
		{
			focus: func() {
				r.tviewApp.SetFocus(r.filter)
			},
		},
	}
}

func (r *TuiMergesRenderer) nextPanel() {
	r.cycleFocus(1)
}

func (r *TuiMergesRenderer) prevPanel() {
	r.cycleFocus(-1)
}

// Or with helper methods
func (r *TuiMergesRenderer) cycleFocus(delta int) {
	N := len(r.focusOps)
	r.currentFocus = ((r.currentFocus+delta)%N + N) % N
	r.focusOps[r.currentFocus].focus()
}

// setupKeyHandlers configures keyboard navigation
func (r *TuiMergesRenderer) setupKeyHandlers() {
	r.tablePage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			r.stop()
			return nil
		case tcell.KeyTab:
			r.nextPanel()
			return nil
		case tcell.KeyBacktab:
			r.prevPanel()
			return nil
		}
		return event
	})

	// Table selection handler
	r.tablePanel.SetSelectedFunc(func(row, col int) {
		if row > 0 { // Skip header row
			r.showTableRowDetails(row - 1)
		}
	})

	// Filter change handler
	r.filter.SetChangedFunc(func(filterText string) {
		r.applyFilter(filterText)
	})
}

// showTableRowDetails displays details for the selected table row
func (r *TuiMergesRenderer) showTableRowDetails(row int) {
	// This would need to be implemented by the specific table type
	// For now, just clear details
	//s := slices.Collect(maps.Values(requests))
	//details.ShowDetails(s[0])

	r.details.Clear()
}

// applyFilter applies the filter to the table content
func (r *TuiMergesRenderer) applyFilter(filterText string) {
	// This would need to be implemented to filter the table content
	// Implementation depends on the specific data being displayed
}
