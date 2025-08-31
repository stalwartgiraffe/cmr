package cmd

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson/jlexer"
	"github.com/spf13/cobra"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/kam"
	"github.com/stalwartgiraffe/cmr/restclient"
)

func NewMergeRequestCommand(app App, cfg *CmdConfig, cancel context.CancelFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "mergerequests",
		Short: "run mergerequests",
		Long:  `Run MergeRequest.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if 0 < len(args) {
				return fmt.Errorf("unexpected args %v", args)
			} else {
				return nil
			}
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
	// start := time.Now()
	accessToken, err := loadGitlabAccessToken()

	app.Println("we got accessToken")

	if err != nil {
		utils.Redln(err)
		return
	}

	mrc := NewMergeRequestClient(accessToken)
	app.Println("start updating recentEvents")
	requests, err := mrc.updateRecentMergeRequest(ctx, app, cancel, filepath, route)
	app.Printf("we got events %d", len(requests))
	if err != nil {
		utils.Redln(err)
		return
	}

	//content := newEventContent(events)
	//promptTable(content, cancel)
}

type MergeRequestClient struct {
	client *gitlab.Client
}

func NewMergeRequestClient(accessToken string) *MergeRequestClient {
	const isVerbose = false
	return &MergeRequestClient{
		client: gitlab.NewClientWithParams(
			"https://gitlab.indexexchange.com/",
			"api/v4/",
			accessToken,
			"xlab",
			isVerbose,
		),
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

	// see https://docs.gitlab.com/ee/api/events.html
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
	app restclient.App,
	resp *resty.Response,
) (*[]gitlab.MergeRequestModel, error) {
	_, span := app.StartSpan(ctx, "unmarshalMergeRequestModel")
	defer span.End()

	if resp == nil {
		return nil, restclient.NewFailureResponse("Response object was nil", resp)
	}
	if resp.IsError() {
		return nil, restclient.NewFailureResponse("ResponseBody="+string(resp.Body()), resp)
	}
	if !resp.IsSuccess() {
		return nil, restclient.NewFailureResponse("Response object had failure status", resp)
	}

	body := resp.Body()

	//err := os.WriteFile("body.txt", body, 0644)
	//if err != nil {
	//	panic(err)
	//}

	//lexer := jlexer.Lexer{Data: resp.Body()}
	lexer := jlexer.Lexer{Data: body}
	var em gitlab.MergeRequestModelSlice

	em.UnmarshalEasyJSON(&lexer)
	if lexer.Error() != nil {
		//panic(lexer.Error())
		app.Println(lexer.Error())
		app.Println(string(body))
		return nil, lexer.Error()
	}
	ss := []gitlab.MergeRequestModel(em)
	return &ss, nil
}

func LoadModels(app App, jsonBlob []byte) (
	*[]gitlab.MergeRequestModel,
	error) {
	lexer := jlexer.Lexer{Data: jsonBlob}
	var em gitlab.MergeRequestModelSlice

	em.UnmarshalEasyJSON(&lexer)
	if lexer.Error() != nil {
		//panic(lexer.Error())
		app.Println(lexer.Error())
		app.Println(string(jsonBlob))
		return nil, lexer.Error()
	}
	ss := []gitlab.MergeRequestModel(em)
	return &ss, nil
}

/*

func promptTable(ptc tview.TableContent, cancel context.CancelFunc) {
	fmt.Println("start promptTable")
	app := tview.NewApplication()
	table := tview.NewTable()

	table.SetContent(ptc)
	table.Select(0, 0).
		SetFixed(1, 1).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				app.Stop()
			}
			if key == tcell.KeyEnter {
				table.SetSelectable(true, true)
			}
		}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(true, true)
	})

	fmt.Println("done table set up")
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
	fmt.Println("finish table ")
}

type eventContent struct {
	tview.TableContentReadOnly

	events []gitlab.EventModel
}

func newEventContent(events gitlab.EventMap) *eventContent {
	s := maps.Values(events)
	sort.Slice(s, func(i, j int) bool {
		return s[i].ID > s[j].ID
	})
	return &eventContent{
		events: s,
	}
}

func (c *eventContent) GetCell(row, col int) *tview.TableCell {
	cell := tview.NewTableCell(c.events[row].Column(col))

	if ((row / 2) % 2) == 0 {
		cell.SetBackgroundColor(tcell.Color234)
	} else {
		cell.SetBackgroundColor(tcell.Color16)
	}
	return cell
}

// Return the total number of rows in the table.
func (c *eventContent) GetRowCount() int {
	return len(c.events)
}

// Return the total number of columns in the table.
func (c *eventContent) GetColumnCount() int {
	return gitlab.EventModelColumnCount()
}
*/
