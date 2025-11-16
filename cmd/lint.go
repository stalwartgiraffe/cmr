package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stalwartgiraffe/cmr/internal/lint"
)

// initCmd represents the init command
func NewLintCommand(cfg *CmdConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "lint",
		Short: "run linters",
		Long: `Run project linters.

FIXME: detailed explanation.
FIXME: include examples.
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				return fmt.Errorf("unexpected args %v", args)
			} else {
				return nil
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			if err := lint.RunEach(ctx, cfg.Config); err != nil {
				fmt.Println(err)
			}
		},
	}

}
