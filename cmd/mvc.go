package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stalwartgiraffe/cmr/internal/tui/merges"
)

func NewMVCCommand(app App, cfg *CmdConfig, cancel context.CancelFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "mvc",
		Short: "run mvc",
		Long:  `Run mvc`,
		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				return fmt.Errorf("unexpected args %v", args)
			} else {
				return nil
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			runMVC(app, cancel, cmd)
		},
	}
}

func runMVC(app App, cancel context.CancelFunc, cmd *cobra.Command) {
	repo := merges.NewInMemoryMergesRepository()
	if err := repo.Load(); err != nil {
		panic(err)
	}
	fmt.Println("mvc2 loaded")
	renderer := merges.NewTuiMergesRenderer()

	controller := merges.NewMergesController(
		repo,
		renderer,
	)

	if err := controller.Run(); err != nil {
		panic(err)
	}
}
