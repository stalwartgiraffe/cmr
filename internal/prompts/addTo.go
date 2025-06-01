package prompts

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/go-git/go-git/v5"
	"github.com/rivo/tview"
	"github.com/stalwartgiraffe/cmr/withstack"
)

func addToWorktree(repo *git.Repository, files []string) {
	fmt.Println("vim-go")
}

type toLineFn func(fp string, s *git.FileStatus) string

func ToWorkTreeStatus(fp string, s *git.FileStatus) string {
	return fmt.Sprintf("%c %s", s.Worktree, fp)
}
func ToStagingStatus(fp string, s *git.FileStatus) string {
	return fmt.Sprintf("%c %s", s.Staging, fp)
}

func SelectFiles(statuses git.Status, toLine toLineFn) ([]string, error) {
	// FIXME implement paneled selector
	// that can move files between git status, the worktree and staging index
	//files, err := prompts.AddFilesToWorktree(files)
	app := tview.NewApplication()
	form := tview.NewForm()
	filePaths, lines := toFileLines(statuses, toLine)
	addLinesText(form, "Files", lines)
	form.
		AddButton(cancelLabel, func() {
			app.Stop()
		})

	isOk := false
	form.
		AddButton(okLabel, func() {
			isOk = true
			app.Stop()
		})

	title := fmt.Sprintf("Add  files worktree")
	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil { // block
		return nil, withstack.Errorf("Could not run form AddFilesTo: %w", err)
	}

	if isOk {
		return filePaths, nil
	}

	return nil, nil
}

func addLinesText(form *tview.Form, label string, lines []string) {
	txt := strings.Join(lines, "\n")
	width := lineWidth(lines, 60)
	height := len(lines)
	form.AddTextArea(label, txt, width, height, len(lines), nil)
}

func toFileLines(statuses git.Status, toLine toLineFn) ([]string, []string) {
	filePaths := maps.Keys(statuses)
	sort.Strings(filePaths)
	lines := []string{}
	for _, filePath := range filePaths {
		status := statuses[filePath]
		lines = append(lines, toLine(filePath, status))
	}
	return filePaths, lines
}

func lineWidth(lines []string, w int) int {
	for _, s := range lines {
		n := len(s)
		if w < n {
			w = n
		}
	}
	return w
}
