package tviewwrapper

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FilterPanel represents a filter input component
type FilterPanel interface {
	GetPrimitive() tview.Primitive
	GetFilter() string
	SetFilter(string)
	SetChangedFunc(func(string))
}

// DetailsPanel represents a details display component
type DetailsPanel interface {
	GetPrimitive() tview.Primitive
	ShowDetails(interface{})
	Clear()
}

// ThreePanelScreen manages a layout with filter, table, and details panels
type ThreePanelScreen struct {
	*tview.Flex

	tviewApp *tview.Application
	filter   FilterPanel
	table    *tview.Table
	details  DetailsPanel

	currentPanel int // 0=filter, 1=table, 2=details
	onStop       StopFunc
}

// NewThreePanelScreen creates a new three-panel layout
func NewThreePanelScreen(
	tviewApp *tview.Application,
	filter FilterPanel,
	tableContent tview.TableContent,
	details DetailsPanel,
	onStop StopFunc,
) *ThreePanelScreen {
	screen := &ThreePanelScreen{
		Flex:         tview.NewFlex(),
		tviewApp:     tviewApp,
		filter:       filter,
		details:      details,
		currentPanel: 1, // Start with table focused
		onStop:       onStop,
	}

	screen.table = MakeContentTable(tableContent, onStop)
	screen.setupLayout()
	screen.setupKeyHandlers()

	return screen
}

// setupLayout arranges the three panels
func (s *ThreePanelScreen) setupLayout() {
	// Main horizontal layout
	s.SetDirection(tview.FlexRow)

	// Filter panel at top (fixed height)
	s.AddItem(s.filter.GetPrimitive(), 3, 0, false)

	// Content area with table and details
	contentFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	contentFlex.AddItem(s.table, 0, 2, true)                   // Table takes 2/3
	contentFlex.AddItem(s.details.GetPrimitive(), 0, 1, false) // Details takes 1/3

	s.AddItem(contentFlex, 0, 1, true)
}

// setupKeyHandlers configures keyboard navigation
func (s *ThreePanelScreen) setupKeyHandlers() {
	s.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			s.onStop()
			return nil
		case tcell.KeyTab:
			s.nextPanel()
			return nil
		case tcell.KeyBacktab:
			s.prevPanel()
			return nil
		}
		return event
	})

	// Table selection handler
	s.table.SetSelectedFunc(func(row, col int) {
		if row > 0 { // Skip header row
			s.showTableRowDetails(row - 1)
		}
	})

	// Filter change handler
	s.filter.SetChangedFunc(func(filterText string) {
		s.applyFilter(filterText)
	})
}

// nextPanel switches to the next panel
func (s *ThreePanelScreen) nextPanel() {
	s.currentPanel = (s.currentPanel + 1) % 3
	s.updateFocus()
}

// prevPanel switches to the previous panel
func (s *ThreePanelScreen) prevPanel() {
	s.currentPanel = (s.currentPanel + 2) % 3
	s.updateFocus()
}

// updateFocus sets focus to the current panel
func (s *ThreePanelScreen) updateFocus() {
	switch s.currentPanel {
	case 0:
		s.tviewApp.SetFocus(s.filter.GetPrimitive())
	case 1:
		s.tviewApp.SetFocus(s.table)
	case 2:
		s.tviewApp.SetFocus(s.details.GetPrimitive())
	}
}

// showTableRowDetails displays details for the selected table row
func (s *ThreePanelScreen) showTableRowDetails(row int) {
	// This would need to be implemented by the specific table type
	// For now, just clear details
	s.details.Clear()
}

// applyFilter applies the filter to the table content
func (s *ThreePanelScreen) applyFilter(filterText string) {
	// This would need to be implemented to filter the table content
	// Implementation depends on the specific data being displayed
}

// GetCurrentPanel returns the currently focused panel (0=filter, 1=table, 2=details)
func (s *ThreePanelScreen) GetCurrentPanel() int {
	return s.currentPanel
}

// SetCurrentPanel sets the focused panel
func (s *ThreePanelScreen) SetCurrentPanel(panel int) {
	if panel >= 0 && panel <= 2 {
		s.currentPanel = panel
		s.updateFocus()
	}
}
