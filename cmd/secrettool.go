package cmd

import (
	"fmt"

	"github.com/stalwartgiraffe/cmr/xr"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
func NewSecretToolCommand(cfg *CmdConfig) *cobra.Command {
	const expectedStatus int = 1
	const secretTool = "secret-tool"
	return &cobra.Command{
		Use:   "st",
		Short: "run secret-tool",
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
			fmt.Println(">secret-tool:")
			stargs := []string{"lookup", "pat", "gitlab"}
			out, err := xr.Run(ctx, secretTool, expectedStatus, nil, stargs...)
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Println(out)
		},
	}
}
