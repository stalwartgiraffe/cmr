package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/stalwartgiraffe/cmr/internal/elog"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/gitutil"
	"gopkg.in/yaml.v3"
)

func NewPullCommand(cfg *CmdConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "pull git repos",
		Long:  `pull well known repos using bulk commands`,

		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				return fmt.Errorf("unexpected args %v", args)
			} else {
				return nil
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := elog.New()

			logger.Println("so pully")

			file, err := os.Open("ignore/projects.yaml")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			var projects []gitlab.ProjectModel
			err = yaml.NewDecoder(file).Decode(&projects)
			if err != nil {
				fmt.Println(err)
				return
			}

			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(err)
				return
			}
			elapsed := []time.Duration{}
			// Now you can use the groups variable
			token, err := loadGitlabAccessToken()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, project := range projects {
				begin := time.Now()
				err := Pull(cfg, project, home, token)
				end := time.Now()

				elapsed = append(elapsed, end.Sub(begin))

				if err != nil {
					fmt.Println("ERROR", err)
				}
			}
			for i, project := range projects {
				fmt.Println(project.NameWithNamespace, elapsed[i])
			}
		},
	}
}

func Pull(cfg *CmdConfig, project gitlab.ProjectModel, home string, token string) error {
	path := gitlab.RepoFilePath(
		home,
		cfg.Config.Repos.Root,
		project,
	)

	fmt.Println(path)
	err := os.MkdirAll(path, os.ModeDir|0755)
	if err != nil {
		return err
	}

	return gitutil.Pull(path, project.HTTPURLToRepo, token, os.Stdout)
}
