package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/utils"
)

func NewViewProjectsCommand(cancel context.CancelFunc, cfg *CmdConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "viewprojects",
		Short: "view projects",
		Long:  `View local projects.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				return fmt.Errorf("unexpected args %v", args)
			} else {
				return nil
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := gitlab.ReadProjects()
			if err != nil {
				utils.Redln(err)
				return
			}

			fmt.Println("number of projects ", len(projects))
		},
	}
}
