package gitlab

import (
	"context"
	"sync"

	"github.com/go-resty/resty/v2"
)

func GatherPageCallsUM[RespT any](
	ctx context.Context,
	client *Client,
	logger AppLog,
	initialQueries <-chan UrlQuery,
	unmarshal func(*resty.Response) (*RespT, error),
) <-chan Call[RespT] {
	return GatherPageCallsWithUM[RespT](
		ctx,
		client,
		logger,
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
	client *Client,
	logger AppLog,
	initialQueries <-chan UrlQuery,
	callCap int,
	queryCap int,
	workersCap int,
	errorCap int,
	unmarshal func(*resty.Response) (*RespT, error),
) <-chan Call[RespT] {
	calls := make([]<-chan Call[RespT], 2)
	var queries <-chan UrlQuery
	calls[0], queries = headPageQueriesUM[RespT](
		ctx,
		client,
		logger,
		initialQueries,
		callCap,
		queryCap,
		errorCap,
		unmarshal,
	)
	calls[1] = tailPageCallsUM[RespT](
		ctx,
		client,
		logger,
		queries,
		workersCap,
		errorCap,
		unmarshal,
	)
	return FanIn(calls)
}

func headPageQueriesUM[RespT any](
	ctx context.Context,
	client *Client,
	logger AppLog,
	firstQueries <-chan UrlQuery,
	callCap int,
	queryCap int,
	errorCap int,
	unmarshal func(*resty.Response) (*RespT, error),
) (
	<-chan Call[RespT],
	<-chan UrlQuery,
) {
	calls := make(chan Call[RespT], callCap)
	queries := make(chan UrlQuery, callCap)
	go func() {
		defer close(calls)
		defer close(queries)
		for firstQuery := range firstQueries {
			//logger.Println("head page GET", firstQuery.Path, firstQuery.Params)
			firstVal, firstHeader, err := GetWithUnmarshal[RespT](
				ctx,
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
	client *Client,
	logger AppLog,
	queries <-chan UrlQuery,
	workersCap int,
	errorsCap int,
	unmarshal func(*resty.Response) (*RespT, error),
) <-chan Call[RespT] {
	calls := make(chan Call[RespT], workersCap)
	go func() {
		defer close(calls)
		var workersWg sync.WaitGroup
		workersWg.Add(workersCap)
		for i := 0; i < workersCap; i++ {
			go func() {
				defer workersWg.Done()
				for q := range queries {
					//logger.Println("tail page GET", q.Path, q.Params)
					v, h, err := GetWithUnmarshal[RespT](
						ctx,
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
