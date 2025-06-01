package cmd

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson/jlexer"
	"github.com/spf13/cobra"
	"github.com/stalwartgiraffe/cmr/internal/elog"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/kam"
	"github.com/stalwartgiraffe/cmr/restclient"
)

func NewMergeRequestCommand(cancel context.CancelFunc, cfg *CmdConfig) *cobra.Command {
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
			filepath := "ignore/my_recent_merge_request.yaml"
			route := "merge_requests/"
			var err error
			// start := time.Now()
			accessToken, err := loadGitlabAccessToken()

			fmt.Println("we got accessToken")

			if err != nil {
				utils.Redln(err)
				return
			}

			mrc := NewMergeRequestClient(accessToken)
			ctx := cmd.Context()
			logger := elog.New()
			fmt.Println("start updating recentEvents")
			requests, err := mrc.updateRecentMergeRequest(ctx, logger, cancel, filepath, route)
			fmt.Printf("we got events %d", len(requests))
			if err != nil {
				utils.Redln(err)
				return
			}

			//content := newEventContent(events)
			//promptTable(content, cancel)
		},
	}
}

/*
func getEvents(
	ctx context.Context,
	logger AppLog,
	cancel context.CancelFunc,
	route string,
	afterThisDate string,
) (
	gitlab.EventMap,
	error,
) {
	var err error
	// start := time.Now()
	accessToken, err := loadGitlabAccessToken()

	if err != nil {
		return nil, err
	}
	ec := NewEventClient(accessToken)
	return ec.getEvents(ctx, logger, cancel, route, afterThisDate)
}
*/

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

		//logger: elog.New(),
	}
}

func (mrc *MergeRequestClient) updateRecentMergeRequest(
	ctx context.Context,
	logger AppLog,
	cancel context.CancelFunc,
	filepath string,
	route string,
) (gitlab.MergeRequestMap, error) {

	//logger.Println("start", route)

	requests, err := gitlab.NewMergeRequestMapFromYaml(filepath)
	if err != nil {
		return nil, err
	}

	recentRequests, err := mrc.getMergeRequests(ctx, logger, cancel, route, requests.LastCreatedDate())
	if err != nil {
		return nil, err
	}
	requests.Insert(recentRequests)
	if err != nil {
		return nil, err
	}
	err = requests.WriteToYamlFile(filepath)
	return requests, nil
}

func (mrc *MergeRequestClient) getMergeRequests(
	ctx context.Context,
	logger AppLog,
	cancel context.CancelFunc,
	route string,
	afterThisDate string,
) (
	gitlab.MergeRequestMap,
	error,
) {
	//logger.Println("getEvents")
	firstQueries := make(chan gitlab.UrlQuery)
	mrCalls := gitlab.GatherPageCallsUM[[]gitlab.MergeRequestModel](
		ctx,
		mrc.client,
		logger,
		firstQueries,
		unmarshalMergeRequestModel,
	)

	//logger.Println("sending first query")
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
	//logger.Println("done sending query")

	requestsMap := gitlab.MergeRequestMap{}
	for s := range mrCalls {
		if s.Error != nil {
			return nil, s.Error
		}
		for _, m := range s.Val {
			requestsMap[m.ID] = m
		}
	}
	//logger.Println("merge calls into map")
	return requestsMap, nil
}

func unmarshalMergeRequestModel(resp *resty.Response) (*[]gitlab.MergeRequestModel, error) {
	if resp == nil {
		return nil, restclient.NewFailureResponse("Response object was nil", resp)
	}
	if resp.IsError() {
		return nil, restclient.NewFailureResponse("ResponseBody="+string(resp.Body()), resp)
	}
	if !resp.IsSuccess() {
		return nil, restclient.NewFailureResponse("Response object had failure status", resp)
	}

	var em gitlab.MergeRequestModelSlice
	body := resp.Body()
	//lexer := jlexer.Lexer{Data: resp.Body()}
	lexer := jlexer.Lexer{Data: body}

	em.UnmarshalEasyJSON(&lexer)
	if lexer.Error() != nil {
		//panic(lexer.Error())
		fmt.Println(lexer.Error())
		fmt.Println(string(body))
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
