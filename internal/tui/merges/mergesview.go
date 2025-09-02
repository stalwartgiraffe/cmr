package merges

import "github.com/rivo/tview"

type MergesRenderer interface {
}

type TuiMergesRenderer struct {
	tviewApp *tview.Application
	stop     StopFn

	tablePage *tview.Flex
}

type StopFn func()

func NewTuiMergesRenderer() {
	r := &TuiMergesRenderer{
		tviewApp: tview.NewApplication(),
	}
	r.stop = r.tviewApp.Stop
}

func (r *TuiMergesRenderer) Run() error {
	return r.tviewApp.SetRoot(r.tablePage, true).SetFocus(r.tablePage).Run()
}
