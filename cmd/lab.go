package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/TwiN/go-color"
	"github.com/spf13/cobra"

	"github.com/stalwartgiraffe/cmr/internal/elog"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/internal/xr"
)

// NewLabCommand initializes the command.
func NewLabCommand(app App, cfg *CmdConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "lab",
		Short: "fetch the collection of projects and write them to projects file",
		Long:  `Run Lab.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				return fmt.Errorf("unexpected args %v", args)
			} else {
				return nil
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			RunLab(app, cmd)
		},
	}
}

func RunLab(app App, cmd *cobra.Command) {
	ctx := cmd.Context()
	ctx, span := app.StartSpan(ctx, "RunLab")
	defer span.End()

	var err error
	accessToken, err := loadGitlabAccessToken()
	if err != nil {
		utils.Redln(err)
		return
	}

	isVerbose := true
	client := gitlab.NewClientWithParams(
		"https://gitlab.indexexchange.com/",
		"api/v4/",
		accessToken,
		"xlab",
		isVerbose,
	)

	const startPage = 1

	logger := elog.New()

	firstQueries := make(chan gitlab.UrlQuery)
	totalPagesLimit := 1000
	//totalPagesLimit := 1
	projectCalls, gatherProjectErrs := gitlab.GatherPageCallsDualApp[[]gitlab.ProjectModel](
		ctx,
		app,
		client,
		logger,
		firstQueries,
		totalPagesLimit,
	)
	errorsFan := []<-chan error{}
	errorsFan = append(errorsFan, gatherProjectErrs)

	firstQueries <- *gitlab.NewPageQuery(
		"projects/",
		startPage,
	)
	close(firstQueries)
	transformCap := 5
	projectResults := gitlab.TransformToOne(
		projectCalls,
		transformCap,
		func(c gitlab.CallNoError[[]gitlab.ProjectModel]) []gitlab.ProjectModel {
			return c.Val
		})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for e := range gitlab.FanIn(errorsFan) {
			fmt.Println(color.Ize(color.Red, e.Error()))
		}
	}()
	projectsMap := make(map[int]gitlab.ProjectModel)
	go func() {
		defer wg.Done()
		for p := range projectResults {
			projectsMap[p.ID] = p
		}
	}()

	wg.Wait()

	fmt.Println("num projects", len(projectsMap))
	if err := utils.WriteToYamlFile("ignore/projects.yaml", utils.ToSortedSlice(projectsMap)); err != nil {
		utils.Redln(err)
		return
	}
	fmt.Println("done reading")
}

func loadGitlabAccessToken() (string, error) {
	token := os.Getenv("GIT_LAB_ACCESS_TOKEN")
	if token != "" {
		return token, nil
	}
	fmt.Println("Xsecret-tool:")
	stArgs := []string{"lookup", "pat", "gitlab"}
	const secretTool = "secret-tool"
	const expectedStatus int = 3
	token, err := xr.Run(secretTool, expectedStatus, nil, stArgs...)
	fmt.Println("pat gitlab", token)
	if err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}

	stArgs = []string{"lookup", "pat", "publicus"}
	token, err = xr.Run(secretTool, expectedStatus, nil, stArgs...)
	fmt.Println("pat publicus", token)
	if err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}
	return "", fmt.Errorf("error: empty accessToken")
}
