package tviewwrapper

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

func TestNewFocusRing(t *testing.T) {
	app := tview.NewApplication()
	panel1 := tview.NewBox()
	panel2 := tview.NewBox()
	panel3 := tview.NewBox()

	tests := []struct {
		name           string
		panels         []tview.Primitive
		expectedPanels int
		expectedFocus  int
	}{
		{
			name:           "empty panels",
			panels:         []tview.Primitive{},
			expectedPanels: 0,
			expectedFocus:  0,
		},
		{
			name:           "single panel",
			panels:         []tview.Primitive{panel1},
			expectedPanels: 1,
			expectedFocus:  0,
		},
		{
			name:           "multiple panels",
			panels:         []tview.Primitive{panel1, panel2, panel3},
			expectedPanels: 3,
			expectedFocus:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ring := NewFocusRing(app, tt.panels...)
			
			require.NotNil(t, ring.tviewApp)
			require.Equal(t, app, ring.tviewApp)
			require.Equal(t, tt.expectedFocus, ring.focusedPanel)
			require.Equal(t, tt.expectedPanels, len(ring.panels))
			require.Equal(t, tt.panels, ring.panels)
		})
	}
}

func TestFocusRing_Cycle(t *testing.T) {
	app := tview.NewApplication()
	panel1 := tview.NewBox()
	panel2 := tview.NewBox()
	panel3 := tview.NewBox()

	tests := []struct {
		name           string
		panels         []tview.Primitive
		initialFocus   int
		direction      RingDirection
		expectedFocus  int
	}{
		{
			name:          "next direction with 3 panels from start",
			panels:        []tview.Primitive{panel1, panel2, panel3},
			initialFocus:  0,
			direction:     NextDir,
			expectedFocus: 1,
		},
		{
			name:          "next direction with 3 panels from middle",
			panels:        []tview.Primitive{panel1, panel2, panel3},
			initialFocus:  1,
			direction:     NextDir,
			expectedFocus: 2,
		},
		{
			name:          "next direction with 3 panels from end (wraps around)",
			panels:        []tview.Primitive{panel1, panel2, panel3},
			initialFocus:  2,
			direction:     NextDir,
			expectedFocus: 0,
		},
		{
			name:          "prev direction with 3 panels from start (wraps around)",
			panels:        []tview.Primitive{panel1, panel2, panel3},
			initialFocus:  0,
			direction:     PrevDir,
			expectedFocus: 2,
		},
		{
			name:          "prev direction with 3 panels from middle",
			panels:        []tview.Primitive{panel1, panel2, panel3},
			initialFocus:  1,
			direction:     PrevDir,
			expectedFocus: 0,
		},
		{
			name:          "prev direction with 3 panels from end",
			panels:        []tview.Primitive{panel1, panel2, panel3},
			initialFocus:  2,
			direction:     PrevDir,
			expectedFocus: 1,
		},
		{
			name:          "single panel next direction",
			panels:        []tview.Primitive{panel1},
			initialFocus:  0,
			direction:     NextDir,
			expectedFocus: 0,
		},
		{
			name:          "single panel prev direction",
			panels:        []tview.Primitive{panel1},
			initialFocus:  0,
			direction:     PrevDir,
			expectedFocus: 0,
		},
		{
			name:          "two panels next from first",
			panels:        []tview.Primitive{panel1, panel2},
			initialFocus:  0,
			direction:     NextDir,
			expectedFocus: 1,
		},
		{
			name:          "two panels next from second (wraps)",
			panels:        []tview.Primitive{panel1, panel2},
			initialFocus:  1,
			direction:     NextDir,
			expectedFocus: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ring := NewFocusRing(app, tt.panels...)
			ring.focusedPanel = tt.initialFocus
			
			ring.Cycle(tt.direction)
			
			require.Equal(t, tt.expectedFocus, ring.focusedPanel)
		})
	}
}

func TestRingDirection_Constants(t *testing.T) {
	tests := []struct {
		name      string
		direction RingDirection
		expected  int
	}{
		{
			name:      "PrevDir constant",
			direction: PrevDir,
			expected:  -1,
		},
		{
			name:      "NextDir constant",
			direction: NextDir,
			expected:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, int(tt.direction))
		})
	}
}

func TestFocusParams(t *testing.T) {
	app := tview.NewApplication()
	panel := tview.NewBox()
	
	params := FocusParams{
		TviewApp: app,
		Panel:    panel,
	}
	
	require.Equal(t, app, params.TviewApp)
	require.Equal(t, panel, params.Panel)
}