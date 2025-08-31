package gitlab

import (
	"context"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/stalwartgiraffe/cmr/restclient"
)

func GatherPageCallsUM[RespT any](
	ctx context.Context,
	app App,
	client *Client,
	initialQueries <-chan UrlQuery,
	unmarshal func(context.Context, restclient.App, *resty.Response) (*RespT, error),
) <-chan Call[RespT] {
	return GatherPageCallsWithUM[RespT](
		ctx,
		app,
		client,
		initialQueries,
		5, // callCap int,
		5, // queryCap int,
		5, // workersCap int,
		1, // errorCap int,
		unmarshal,
	)
}
func GatherPageCallsWithUM[RespT any](
	ctx context.Context,
	app App,
	client *Client,
	initialQueries <-chan UrlQuery,
	callCap int,
	queryCap int,
	workersCap int,
	errorCap int,
	unmarshal func(context.Context, restclient.App, *resty.Response) (*RespT, error),
) <-chan Call[RespT] {
	ctx, span := app.StartSpan(ctx, "GatherPageCallsWithUM")
	defer span.End()

	calls := make([]<-chan Call[RespT], 2)
	var queries <-chan UrlQuery
	calls[0], queries = headPageQueriesUM[RespT](
		ctx,
		app,
		client,
		initialQueries,
		callCap,
		queryCap,
		errorCap,
		unmarshal,
	)
	calls[1] = tailPageCallsUM[RespT](
		ctx,
		app,
		client,
		queries,
		workersCap,
		errorCap,
		unmarshal,
	)
	return FanIn(calls)
}

func headPageQueriesUM[RespT any](
	ctx context.Context,
	app App,
	client *Client,
	firstQueries <-chan UrlQuery,
	callCap int,
	queryCap int,
	errorCap int,
	unmarshal func(context.Context, restclient.App, *resty.Response) (*RespT, error),
) (
	<-chan Call[RespT],
	<-chan UrlQuery,
) {
	ctx, span := app.StartSpan(ctx, "headPageQueriesUM")
	defer span.End()

	calls := make(chan Call[RespT], callCap)
	queries := make(chan UrlQuery, callCap)
	go func() {
		ctx, span := app.StartSpan(ctx, "go_headPageQueriesUM")
		defer span.End()
		defer close(calls)
		defer close(queries)
		for firstQuery := range firstQueries {
			firstVal, firstHeader, err := GetWithUnmarshal[RespT](
				ctx,
				app,
				client,
				firstQuery.Path,
				firstQuery.Params,
				unmarshal,
			)
			if err != nil {
				calls <- Call[RespT]{
					Query: firstQuery,
					Error: err,
				}
				continue
			} else {
				calls <- Call[RespT]{
					Query:  firstQuery,
					Header: firstHeader,
					Val:    *firstVal,
				}
			}

			cursor, err := parsePageCursor(firstHeader)
			if err != nil {
				calls <- Call[RespT]{
					Query: firstQuery,
					Error: err,
				}
				continue
			}
			if cursor.page == nil || cursor.totalPages == nil {
				continue
			}
			p := *cursor.page + 1
			n := *cursor.totalPages
			for ; p <= n; p++ {
				next := *firstQuery.Clone()
				next.Params["page"] = p
				queries <- next
			}
		}
	}()
	return calls, queries
}

func tailPageCallsUM[RespT any](
	ctx context.Context,
	app App,
	client *Client,
	queries <-chan UrlQuery,
	workersCap int,
	errorsCap int,
	unmarshal func(context.Context, restclient.App, *resty.Response) (*RespT, error),
) <-chan Call[RespT] {
	ctx, span := app.StartSpan(ctx, "tailPageQueriesUM")
	defer span.End()

	calls := make(chan Call[RespT], workersCap)
	go func() {
		ctx, span := app.StartSpan(ctx, "go_tailPageQueriesUM")
		defer span.End()
		defer close(calls)
		var workersWg sync.WaitGroup
		workersWg.Add(workersCap)
		for i := 0; i < workersCap; i++ {
			go func() {
				ctx, span := app.StartSpan(ctx, "go_go_tailPageQueriesUM")
				defer span.End()

				defer workersWg.Done()
				for q := range queries {
					v, h, err := GetWithUnmarshal[RespT](
						ctx,
						app,
						client,
						q.Path,
						q.Params,
						unmarshal,
					)
					if err != nil {
						calls <- Call[RespT]{
							Query: q,
							Error: err,
						}
						continue
					}

					calls <- Call[RespT]{
						Query:  q,
						Header: h,
						Val:    *v,
					}
				}
			}()
		}
		workersWg.Wait()

	}()
	return calls
}
