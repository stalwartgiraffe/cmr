package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson/jlexer"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"

	"github.com/stalwartgiraffe/cmr/internal/gitlab"
	"github.com/stalwartgiraffe/cmr/internal/tviewwrapper"
	"github.com/stalwartgiraffe/cmr/internal/utils"
	"github.com/stalwartgiraffe/cmr/kam"
	"github.com/stalwartgiraffe/cmr/restclient"
	"github.com/stalwartgiraffe/cmr/withstack"
)

func NewEventsCommand(app App, cfg *CmdConfig, cancel context.CancelFunc) *cobra.Command {
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
			runEventsCmd(app, cancel, cmd)
		},
	}
}

func runEventsCmd(app App, cancel context.CancelFunc, cmd *cobra.Command) {
	ctx := cmd.Context()
	ctx, span := app.StartSpan(ctx, "runEventsCmd")
	defer span.End()

	projects, err := gitlab.ReadProjects()
	if err != nil {
		utils.Redln(err)
		return
	}

	filepath := "ignore/my_recent_events.yaml"
	route := "events/"
	accessToken, err := loadGitlabAccessToken()
	if err != nil {
		utils.Redln(err)
		return
	}

	ec := NewEventClient(accessToken)
	app.Println("start updating recentEvents")
	events, err := ec.updateRecentEvents(ctx, app, cancel, filepath, route)
	if err != nil {
		utils.Redln(err)
		return
	}
	app.Printf("we got events %d", len(events))
	_ = events

	content := tviewwrapper.NewEventsContent(events, projects)
	appTableRun(content, cancel)
}

func getEvents(
	ctx context.Context,
	app App,
	cancel context.CancelFunc,
	route string,
	afterThisDate string,
) (
	gitlab.EventMap,
	error,
) {
	var err error
	accessToken, err := loadGitlabAccessToken()
	if err != nil {
		return nil, err
	}
	ec := NewEventClient(accessToken)
	return ec.getEvents(ctx, app, cancel, route, afterThisDate)
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
	app App,
	cancel context.CancelFunc,
	filepath string,
	route string,
) (gitlab.EventMap, error) {
	ctx, span := app.StartSpan(ctx, "updateRecentEvents")
	defer span.End()

	events, err := gitlab.NewEventMapFromYaml(ctx, app, filepath)
	if err != nil {
		return nil, err
	}

	recentEvents, err := ec.getEvents(ctx, app, cancel, route, events.LastDate())
	if err != nil {
		return nil, err
	}
	events.Insert(recentEvents)
	err = events.WriteToYamlFile(filepath)
	return events, err
}

func unmarshalEventModel(
	ctx context.Context,
	app restclient.App,
	resp *resty.Response,
) (*[]gitlab.EventModel, error) {
	_, span := app.StartSpan(ctx, "unmarshalEventModel")
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

	var em gitlab.EventModelSlice
	lexer := jlexer.Lexer{Data: resp.Body()}
	em.UnmarshalEasyJSON(&lexer)
	if lexer.Error() != nil {
		return nil, lexer.Error()
	}
	ss := []gitlab.EventModel(em)
	return &ss, nil
}

// verifyAllFieldsExpected returns an error if a field is not expected
func verifyAllFieldsExpected(data []byte, names map[string]struct{}) error {
	kvs := []map[string]any{}
	if err := json.Unmarshal(data, &kvs); err != nil {
		return withstack.Errorf("Unmarshal error:%w", err)
	}

	for _, kv := range kvs {
		for k, v := range kv {
			if _, ok := names[k]; !ok {
				return fmt.Errorf("unexpected field name %s %v %d", k, v, data)
			}
		}
	}
	return nil
}

func (ec *EventClient) getEvents(
	ctx context.Context,
	app App,
	_ context.CancelFunc,
	route string,
	afterThisDate string,
) (
	gitlab.EventMap,
	error,
) {
	ctx, span := app.StartSpan(ctx, "getEvents")
	defer span.End()

	firstQueries := make(chan gitlab.UrlQuery)
	eventCalls := gitlab.GatherPageCallsUM[[]gitlab.EventModel](
		ctx,
		app,
		ec.client,
		firstQueries,
		unmarshalEventModel,
	)

	// see https://docs.gitlab.com/ee/api/events.html
	const startPage = 1
	const per_page = 200 // nolint
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
	for s := range eventCalls {
		if s.Error != nil {
			return nil, s.Error
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
