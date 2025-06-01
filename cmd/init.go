package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func addInitCommand(cfg *CmdConfig, parent *cobra.Command) {
	initCmd := NewInitCommand()
	addBranchCommand(initCmd)
	parent.AddCommand(initCmd)
}

// initCmd represents the init command
func NewInitCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize an object",
		Long: `Initialize an object.

FIXME: detailed explanation.
FIXME: include examples.
`,
		Args: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("init must be called with a sub command")
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("init called")
		},
	}

	return initCmd
}
