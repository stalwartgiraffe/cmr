package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson/jlexer"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/stalwartgiraffe/cmr/internal/elog"
	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/tviewwrapper"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/kam"
	"github.com/stalwartgiraffe/cmr/restclient"
	"github.com/stalwartgiraffe/cmr/withstack"
)

type AppLog interface {
	Printf(format string, v ...any)
	Print(v ...any)
	Println(v ...any)
	Flush()
}

func NewEventsCommand(cancel context.CancelFunc, cfg *CmdConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "events",
		Short: "run events",
		Long:  `Run Events.`,
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

			filepath := "ignore/my_recent_events.yaml"
			route := "events/"
			// start := time.Now()
			accessToken, err := loadGitlabAccessToken()

			fmt.Println("we got accessToken")

			if err != nil {
				utils.Redln(err)
				return
			}

			ec := NewEventClient(accessToken)
			ctx := cmd.Context()
			logger := elog.New()
			fmt.Println("start updating recentEvents")
			events, err := ec.updateRecentEvents(ctx, logger, cancel, filepath, route)
			fmt.Printf("we got events %d", len(events))
			_ = events
			if err != nil {
				utils.Redln(err)
				return
			}

			content := tviewwrapper.NewEventsContent(events, projects)
			appTableRun(content, cancel)
		},
	}
}

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

type EventClient struct {
	client *gitlab.Client
}

func NewEventClient(accessToken string) *EventClient {
	const isVerbose = false
	return &EventClient{
		client: gitlab.NewClientWithParams(
			"https://gitlab.indexexchange.com/",
			"api/v4/",
			accessToken,
			"xlab",
			isVerbose,
		),
	}
}

func (ec *EventClient) updateRecentEvents(
	ctx context.Context,
	logger AppLog,
	cancel context.CancelFunc,
	filepath string,
	route string,
) (gitlab.EventMap, error) {
	events, err := gitlab.NewEventMapFromYaml(filepath)
	if err != nil {
		return nil, err
	}

	recentEvents, err := ec.getEvents(ctx, logger, cancel, route, events.LastDate())
	if err != nil {
		return nil, err
	}
	events.Insert(recentEvents)
	err = events.WriteToYamlFile(filepath)
	return events, err
}

func unmarshalEventModel(resp *resty.Response) (*[]gitlab.EventModel, error) {
	if resp == nil {
		return nil, restclient.NewFailureResponse("Response object was nil", resp)
	}
	if resp.IsError() {
		return nil, restclient.NewFailureResponse("ResponseBody="+string(resp.Body()), resp)
	}
	if !resp.IsSuccess() {
		return nil, restclient.NewFailureResponse("Response object had failure status", resp)
	}

	var em gitlab.EventModelSlice
	lexer := jlexer.Lexer{Data: resp.Body()}
	em.UnmarshalEasyJSON(&lexer)
	if lexer.Error() != nil {
		return nil, lexer.Error()
	}
	ss := []gitlab.EventModel(em)
	return &ss, nil
}

// veryifyAllFieldsExpected returns an error if a field is not expected
func veryifyAllFieldsExpected(data []byte, names map[string]struct{}) error {
	kvs := []map[string]any{}
	if err := json.Unmarshal(data, &kvs); err != nil {
		return withstack.Errorf("Unmarshal error:%w", err)
	}

	for _, kv := range kvs {
		for k, v := range kv {
			if _, ok := names[k]; !ok {
				return fmt.Errorf("unexpectd field name %s %v %d", k, v, data)
			}
		}
	}
	return nil
}

func (ec *EventClient) getEvents(
	ctx context.Context,
	logger AppLog,
	_ context.CancelFunc,
	route string,
	afterThisDate string,
) (
	gitlab.EventMap,
	error,
) {
	firstQueries := make(chan gitlab.UrlQuery)
	eventCalls := gitlab.GatherPageCallsUM[[]gitlab.EventModel](
		ctx,
		ec.client,
		logger,
		firstQueries,
		unmarshalEventModel,
	)

	// see https://docs.gitlab.com/ee/api/events.html
	const startPage = 1
	const per_page = 200
	firstQueries <- gitlab.UrlQuery{
		Path: route,
		Params: kam.Map{
			// action - include only particular action type
			// target_type - include only a particular target type
			"after":    afterThisDate, // 2006-01-02 format expected
			"sort":     "desc",        // newest first
			"page":     startPage,
			"per_page": per_page,
		},
	}

	close(firstQueries)

	eventsMap := gitlab.EventMap{}
	var err error
	for s := range eventCalls {
		if s.Error != nil {
			return nil, err
		}
		for _, m := range s.Val {
			eventsMap[m.ID] = m
		}
	}

	return eventsMap, nil
}

func appTableRun(ptc tview.TableContent, _ context.CancelFunc) {
	app := tview.NewApplication()
	table := tviewwrapper.MakeContentTable(ptc, app.Stop)
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
