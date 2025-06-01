package cmd

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
func NewPushCommand(cfg *CmdConfig, repo *git.Repository) *cobra.Command {
	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "run pushers",
		Long: `Run project pushers.

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
			RunPush()
		},
	}

	return pushCmd
}
func RunPush() {
}
