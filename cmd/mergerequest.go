package cmd

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson/jlexer"
	"github.com/spf13/cobra"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	tw "github.com/stalwartgiraffe/cmr/internal/tviewwrapper"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/kam"
	rc "github.com/stalwartgiraffe/cmr/restclient"
)

func NewMergeRequestCommand(app App, cfg *CmdConfig, cancel context.CancelFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "mergerequests",
		Short: "run mergerequests",
		Long:  `Run MergeRequest.`,
		Args: func(cmd *cobra.Command, args []string) error {
			return NoArgs(args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runMergeRequestCmd(app, cancel, cmd)
		},
	}
}
func runMergeRequestCmd(app App, cancel context.CancelFunc, cmd *cobra.Command) {
	ctx := cmd.Context()
	ctx, span := app.StartSpan(ctx, "runMergeRequestCmd")
	defer span.End()

	filepath := "ignore/my_recent_merge_request.yaml"
	route := "merge_requests/"
	var err error
	authToken, err := loadGitlabAuthToken()

	app.Println("we got authToken")

	if err != nil {
		utils.Redln(err)
		return
	}

	client := NewMergeRequestClient(
		rc.WithBaseURL("https://gitlab.indexexchange.com/"),
		rc.WithAuthToken(authToken),
		rc.WithIsVerbose(true),
	)
	app.Println("start updating recentEvents")
	requests, err := client.updateRecentMergeRequest(ctx, app, cancel, filepath, route)
	app.Printf("we got events %d", len(requests))
	if err != nil {
		utils.Redln(err)
		return
	}

	tableContent := tw.NewTwoBandTableContent(
		tw.NewMergeRequestTextTable(
			nil, requests,
		),
	)

	appTableRun(tableContent, cancel)
}

type MergeRequestClient struct {
	client *gitlab.Client
}

func NewMergeRequestClient(overrides ...rc.Option) *MergeRequestClient {
	return &MergeRequestClient{
		client: gitlab.NewClient(overrides...),
	}
}

func (mrc *MergeRequestClient) updateRecentMergeRequest(
	ctx context.Context,
	app App,
	cancel context.CancelFunc,
	filepath string,
	route string,
) (gitlab.MergeRequestMap, error) {
	ctx, span := app.StartSpan(ctx, "updateRecentMergeRequest")
	defer span.End()

	requests, err := gitlab.NewMergeRequestMapFromYaml(filepath)
	if err != nil {
		return nil, err
	}

	recentRequests, err := mrc.getMergeRequests(ctx, app, cancel, route, requests.LastCreatedDate())
	if err != nil {
		return nil, err
	}
	requests.Insert(recentRequests)
	err = requests.WriteToYamlFile(filepath)
	return requests, err
}

func (mrc *MergeRequestClient) getMergeRequests(
	ctx context.Context,
	app App,
	cancel context.CancelFunc,
	route string,
	afterThisDate string,
) (
	gitlab.MergeRequestMap,
	error,
) {
	ctx, span := app.StartSpan(ctx, "getMergeRequests")
	defer span.End()

	firstQueries := make(chan gitlab.UrlQuery)
	mrCalls := gitlab.GatherPageCallsUM[[]gitlab.MergeRequestModel](
		ctx,
		app,
		mrc.client,
		firstQueries,
		unmarshalMergeRequestModel,
	)

	// see https://docs.gitlab.com/api/merge_requests/
	const startPage = 1
	const per_page = 200
	firstQueries <- gitlab.UrlQuery{
		Path: route,
		Params: kam.Map{
			// action - include only particular action type
			// target_type - include only a particular target type

			"updated_after": afterThisDate, // 2006-01-02 format expected
			"sort":          "desc",        // newest first

			"page":     startPage,
			"per_page": per_page,
		},
	}

	close(firstQueries)

	requestsMap := gitlab.MergeRequestMap{}
	for s := range mrCalls {
		if s.Error != nil {
			return nil, s.Error
		}
		for _, m := range s.Val {
			requestsMap[m.ID] = m
		}
	}
	return requestsMap, nil
}

func unmarshalMergeRequestModel(
	ctx context.Context,
	app rc.App,
	resp *resty.Response,
) (*[]gitlab.MergeRequestModel, error) {
	_, span := app.StartSpan(ctx, "unmarshalMergeRequestModel")
	defer span.End()

	if resp == nil {
		return nil, rc.NewFailureResponse("Response object was nil", resp)
	}
	if resp.IsError() {
		return nil, rc.NewFailureResponse("ResponseBody="+string(resp.Body()), resp)
	}
	if !resp.IsSuccess() {
		return nil, rc.NewFailureResponse("Response object had failure status", resp)
	}

	return unmarshalModels(app, resp.Body())
}

var countModel int

func unmarshalModels(app App, jsonBlob []byte) (
	*[]gitlab.MergeRequestModel,
	error) {

	prettyJson, err := utils.PrettyJSON(jsonBlob)
	if err == nil {
		exampleName := fmt.Sprintf("internal/gitlab/localhost/api/merge_request_examples/merge_request_%03d.json", countModel)
		countModel++
		utils.WriteStringToFile(exampleName, prettyJson)
	}

	lexer := jlexer.Lexer{Data: jsonBlob}
	var em gitlab.MergeRequestModelSlice

	em.UnmarshalEasyJSON(&lexer)
	if lexer.Error() != nil {
		errTxt := lexer.Error().Error()
		body := string(jsonBlob)
		prettyTxt := prettySubStringJson(body, errTxt)
		return nil, fmt.Errorf("Error:%s\n%s", errTxt, prettyTxt)
	}

	ss := []gitlab.MergeRequestModel(em)
	return &ss, nil
}

func prettySubStringJson(body string, errTxt string) string {
	last := len(body) - 1
	if last < 0 {
		return ""
	}

	mid, err := parseNextInt(errTxt)
	if err != nil {
		return body
	}

	start := lastNIndex(body[:mid-1], 5, ",")
	start = max(0, start)
	end := nextNIndex(body[mid+1:], 5, ",")
	if end == -1 {
		end = last
	} else {
		end = min(last, mid+end+1)
	}

	txt := body[start:end]
	if strings.Contains(txt, "\n") {
		return txt
	} else {
		return strings.Replace(txt, ",", ",\n", -1)
	}
}

func subStringJson(body string, mid int) string {
	start := lastNIndex(body[:mid-1], 10, ",")
	start = max(0, start)
	end := lastNIndex(body[mid+1:], 10, ",")
	end = min(len(body)-1, end)

	txt := body[start:end]
	return strings.Replace(txt, ",", ",\n", -1)
}

var reInt = regexp.MustCompile(`\d+`)

func parseNextInt(s string) (int, error) {
	match := reInt.FindString(s)
	if match == "" {
		return 0, fmt.Errorf("no number found in string")
	}
	return strconv.Atoi(match)
}

func lastNIndex(body string, count int, txt string) int {
	currentPos := len(body)
	for range count {
		last := strings.LastIndex(body[:currentPos], txt)
		if last == -1 {
			return -1
		}
		currentPos = last
	}
	return currentPos
}

func nextNIndex(body string, count int, txt string) int {
	if count == 0 {
		return -1
	}
	currentPos := 0
	found := 0
	for range count {
		found = strings.Index(body[currentPos:], txt)
		if found == -1 {
			return -1
		}
		currentPos = currentPos + found + 1
	}
	return currentPos
}
