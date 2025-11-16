package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	rc "github.com/stalwartgiraffe/cmr/restclient"
	"github.com/stalwartgiraffe/cmr/xr"
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
			runLabCmd(app, cmd)
		},
	}
}

func runLabCmd(app App, cmd *cobra.Command) {
	ctx := cmd.Context()
	ctx, span := app.StartSpan(ctx, "RunLab")
	defer span.End()

	var err error
	authToken, err := loadGitlabAuthToken(ctx)
	if err != nil {
		utils.Redln(err)
		return
	}

	client := NewProjectsClient(
		rc.WithBaseURL("https://gitlab.indexexchange.com/"),
		rc.WithAuthToken(authToken),
	)
	projects, errs := client.getProjects(
		ctx,
		app)

	if errs != nil {
		utils.Redln(errs)
		return
	}

	fmt.Println("num projects", len(projects))
	if err := utils.WriteToYamlFile("ignore/projects.yaml", utils.ToSortedSlice(projects)); err != nil {
		utils.Redln(err)
		return
	}
	fmt.Println("done reading")
}

type ProjectsClient struct {
	client *gitlab.Client
}

func NewProjectsClient(overrides ...rc.Option) *ProjectsClient {
	return &ProjectsClient{
		client: gitlab.NewClient(overrides...),
	}
}

func (pc *ProjectsClient) getProjects(
	ctx context.Context,
	app App,
) (map[int]gitlab.ProjectModel, error) {
	const startPage = 1

	firstQueries := make(chan gitlab.UrlQuery)
	totalPagesLimit := 1000
	projectCalls, gatherProjectErrs := gitlab.GatherPageCallsDualApp[[]gitlab.ProjectModel](
		ctx,
		app,
		pc.client,
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
	var errs error
	go func() {
		defer wg.Done()
		for e := range gitlab.FanIn(errorsFan) {
			errs = errors.Join(errs, e)
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
	return projectsMap, errs
}

func loadGitlabAuthToken(ctx context.Context) (string, error) {
	token := os.Getenv("GIT_LAB_ACCESS_TOKEN")
	if token != "" {
		return token, nil
	}
	fmt.Println("Xsecret-tool:")
	stArgs := []string{"lookup", "pat", "gitlab"}
	const secretTool = "secret-tool"
	const expectedStatus int = 3
	token, err := xr.Run(ctx, secretTool, expectedStatus, nil, stArgs...)
	fmt.Println("pat gitlab", token)
	if err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}

	stArgs = []string{"lookup", "pat", "publicus"}
	token, err = xr.Run(ctx, secretTool, expectedStatus, nil, stArgs...)
	fmt.Println("pat publicus", token)
	if err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}
	return "", fmt.Errorf("error: empty authToken")
}
