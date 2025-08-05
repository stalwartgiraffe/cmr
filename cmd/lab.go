package cmd

import (
	//"atomic"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/TwiN/go-color"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"

	"github.com/stalwartgiraffe/cmr/internal/elog"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/internal/xr"
	"github.com/stalwartgiraffe/cmr/kam"
)

// NewLabCommand initalizes the command.
func NewLabCommand(app App, cfg *CmdConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "lab",
		Short: "run lab",
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
	if app != nil {
		var span trace.Span
		ctx, span = app.StartSpan(ctx, "RunLab")
		defer span.End()
	}

	var err error
	//start := time.Now()
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
	//loginTime := time.Now()

	const startPage = 1

	logger := elog.New()

	firstQueries := make(chan gitlab.UrlQuery)
	totalPagesLimit := 1
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
	//var count atomic.Int64 // <-- added atomic counter
	projectResults := gitlab.TransformToOne(
		projectCalls,
		transformCap,
		func(c gitlab.CallNoError[[]gitlab.ProjectModel]) []gitlab.ProjectModel {
			return c.Val
			/*
				result := []gitlab.ProjectModel{}
				for _, project := range c.Val {

					fmt.Println(project.ID)
					result = append(result, project)
					// groups[group.ID] = group
					// queries = append(queries,
					// 	*gitlab.NewPageQuery(
					// 		fmt.Sprintf("groups/%d/projects", group.ID),
					// 		startPage,
					// 	))
				}
				return result
			*/
		})
	// fmt.Println("count of project", len(projects))
	// for e := range gitlab.FanIn(errorsFan) {
	// 	fmt.Println(color.Ize(color.Red, e.Error()))
	// }
	//
	// //projects := make(map[int]gitlab.ProjectModel)
	// //for projects := range projectCalls {
	// //fmt.Println(len(projects), "--------------------------------")
	// //for _, p := range projects {
	// //fmt.Println(p.ID)
	// //}
	// //}

	// setupTime := time.Now()
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
			// for _, p := range result {
			projectsMap[p.ID] = p
			// }
		}
	}()

	wg.Wait()
	// readTime := time.Now()

	// fmt.Println("num groups", len(groups))
	fmt.Println("num projects", len(projectsMap))
	if err := utils.WriteToYamlFile("ignore/projects.yaml", utils.ToSortedSlice(projectsMap)); err != nil {
		utils.Redln(err)
		return
	}
	fmt.Println("done reading")
}

func oldRun(cmd *cobra.Command) {
	ctx := cmd.Context()
	var err error
	start := time.Now()
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
	loginTime := time.Now()

	const startPage = 1

	logger := elog.New()

	firstQueries := make(chan gitlab.UrlQuery)
	totalPagesLimit := 1
	groupsCalls, gatherGroupErrs := gitlab.GatherPageCallsDual[[]gitlab.GroupModel](
		ctx,
		client,
		logger,
		firstQueries,
		totalPagesLimit,
	)
	errorsFan := []<-chan error{}
	errorsFan = append(errorsFan, gatherGroupErrs)

	firstQueries <- *gitlab.NewPageQuery(
		"groups/",
		startPage,
	)
	close(firstQueries)
	groups := make(map[int]gitlab.GroupModel)
	transformCap := 5
	projectQueries := gitlab.TransformToOne(
		groupsCalls,
		transformCap,
		func(c gitlab.CallNoError[[]gitlab.GroupModel]) []gitlab.UrlQuery {
			queries := []gitlab.UrlQuery{}
			for _, group := range c.Val {
				groups[group.ID] = group
				queries = append(queries,
					*gitlab.NewPageQuery(
						fmt.Sprintf("groups/%d/projects", group.ID),
						startPage,
					))
			}
			return queries
		})

	projectCalls, projectErrs := gitlab.GatherPageCallsDual[[]gitlab.ProjectModel](
		ctx,
		client,
		logger,
		projectQueries,
		totalPagesLimit,
	)

	errorsFan = append(errorsFan, projectErrs)

	setupTime := time.Now()
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
		for c := range projectCalls {

			for _, p := range c.Val {
				projectsMap[p.ID] = p
			}
		}
	}()

	wg.Wait()
	readTime := time.Now()

	fmt.Println("num groups", len(groups))
	fmt.Println("num sgIDs", len(projectsMap))
	fmt.Println("done reading")
	utils.WriteToYamlFile("ignore/groups.yaml", utils.ToSortedSlice(groups))
	utils.WriteToYamlFile("ignore/projects.yaml", utils.ToSortedSlice(projectsMap))
	fmt.Println("done writing")
	writeTime := time.Now()

	fmt.Println("Elapsed",
		"login", loginTime.Sub(start),
		"setup", setupTime.Sub(loginTime),
		"read", readTime.Sub(setupTime),
		"write", writeTime.Sub(readTime),
	)
}

func dumproute(
	ctx context.Context,
	client *gitlab.Client) {
	q := gitlab.UrlQuery{
		Path: "groups/45/descendant_groups",

		Params: kam.Map{
			"order_by":               "id",
			"owned":                  false,
			"page":                   1,
			"per_page":               200,
			"sort":                   "asc",
			"statistics":             false,
			"with_custom_attributes": false,
		},
	}
	v, m, e := client.Get(ctx, q)
	if e != nil {
		fmt.Println("ERROR:", e)
	} else {
		_ = m
		fmt.Println("- descix ------------------------------------------")
		fmt.Println(utils.YamlString(v))

		_ = e
	}

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
