package cmd

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/tviewwrapper"
	"github.com/stalwartgiraffe/cmr/internal/utils"
)

func NewDemoCommand(app App, cfg *CmdConfig, cancel context.CancelFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "demo",
		Short: "run demo",
		Long:  `Run demo`,
		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				return fmt.Errorf("unexpected args %v", args)
			} else {
				return nil
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			runDemo(app, cancel, cmd)
		},
	}
}
func runDemo(app App, cancel context.CancelFunc, cmd *cobra.Command) {
	projects, err := gitlab.ReadProjects()
	if err != nil {
		utils.Redln(err)
		return
	}
	filepath := "ignore/my_recent_merge_request.yaml"
	requests, err := gitlab.NewMergeRequestMapFromYaml(filepath)
	if err != nil {
		utils.Redln(err)
		return
	}

	tableContent := tviewwrapper.NewTwoBandTableContent(
		tviewwrapper.NewMergeRequestTextTable(
			projects,
			requests,
		))

	tviewApp := tview.NewApplication()
	filter := tviewwrapper.NewBasicFilter("sure")
	details := tviewwrapper.NewTextDetails()

	s := slices.Collect(maps.Values(requests))
	details.ShowDetails(s[0])
	screen := tviewwrapper.NewThreePanelScreen(
		tviewApp,
		filter,
		tableContent,
		details,
		tviewApp.Stop,
	)

	if err := tviewApp.SetRoot(screen, true).SetFocus(screen).Run(); err != nil {
		panic(err)
	}
}
