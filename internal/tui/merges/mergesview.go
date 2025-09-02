package merges

import (
	"github.com/rivo/tview"

	"github.com/stalwartgiraffe/cmr/internal/tviewwrapper"
)

type MergesRenderer interface {
	Run() error

	MakeBinding(repo MergesRepository)
}

type TuiMergesRenderer struct {
	tviewApp *tview.Application
	stop     StopFn

	tablePage *tview.Flex

	textTable  tviewwrapper.TextTable
	tablePanel *tview.Table
}

type StopFn func()

func NewTuiMergesRenderer() *TuiMergesRenderer {
	r := &TuiMergesRenderer{
		tviewApp:  tview.NewApplication(),
		tablePage: tview.NewFlex(),
	}
	r.stop = r.tviewApp.Stop

	return r
}

func (r *TuiMergesRenderer) Run() error {
	return r.tviewApp.SetRoot(r.tablePage, true).SetFocus(r.tablePage).Run()
}

func (r *TuiMergesRenderer) MakeBinding(repo MergesRepository) {
	projects, mergesMap := repo.GetCollections()
	r.textTable = tviewwrapper.NewMergeRequestTextTable(
		projects,
		mergesMap,
	)
	tableContent := tviewwrapper.NewTwoBandTable(r.textTable)
	r.tablePanel = tviewwrapper.MakeContentTable(tableContent, r.tviewApp.Stop)

	r.setupLayout()
}

func (r *TuiMergesRenderer) setupLayout() {
	// Main horizontal layout
	r.tablePage.SetDirection(tview.FlexRow)
	r.tablePage.AddItem(r.tablePanel, 0, 1, true)
}
