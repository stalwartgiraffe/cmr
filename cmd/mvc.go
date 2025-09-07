package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stalwartgiraffe/cmr/internal/reload"
	"github.com/stalwartgiraffe/cmr/internal/tui/merges"
)

func NewMVCCommand(app App, cfg *CmdConfig, cancel context.CancelFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "mvc",
		Short: "run mvc",
		Long:  `Run mvc`,
		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				for _, arg := range args {
					switch arg {
					case "reload":
						if err := reload.BeginWatchPwd(cancel); err != nil {
							return err
						}
					default:
						return fmt.Errorf("Unknown argument %s", arg)
					}
				}
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runMVC(cancel, app, cmd)
		},
	}
}

func NoArgs(args []string) error {
	if 0 < len(args) {
		return fmt.Errorf("unexpected args %v", args)
	} else {
		return nil
	}
}

func runMVC(cancel context.CancelFunc, app App, cmd *cobra.Command) {
	ctx := cmd.Context()
	repo := merges.NewInMemoryMergesRepository()

	// TODO handle in go rouine
	if err := repo.Load(); err != nil {
		panic(err)
	}
	renderer := merges.NewTuiMergesRenderer(ctx, repo)

	controller := merges.NewMergesController(
		repo,
		renderer,
	)

	if err := controller.Run(); err != nil {
		panic(err)
	}
}
