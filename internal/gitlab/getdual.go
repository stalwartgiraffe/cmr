package gitlab

import (
	"context"
	"sync"
)

type AppLog interface {
	Printf(format string, v ...any)
	Print(v ...any)
	Println(v ...any)
}

func GatherPageCallsDual[RespT any](
	ctx context.Context,
	client *Client,
	logger AppLog,
	initialQueries <-chan UrlQuery,
	totalPagesLimit int,
) (
	<-chan CallNoError[RespT],
	<-chan error,
) {
	return GatherPageCallsWithDual[RespT](
		ctx,
		client,
		logger,
		initialQueries,
		5,               // callCap int,
		5,               // queryCap int,
		5,               // workersCap int,
		1,               // errorCap int,
		totalPagesLimit, // 0 means no limit
	)
}

func GatherPageCallsWithDual[RespT any](
	ctx context.Context,
	client *Client,
	logger AppLog,
	initialQueries <-chan UrlQuery,
	callCap int,
	queryCap int,
	workersCap int,
	errorCap int,
	totalPagesLimit int, // 0 means no limit
) (
	<-chan CallNoError[RespT],
	<-chan error,
) {
	calls := make([]<-chan CallNoError[RespT], 2)
	errors := make([]<-chan error, 2)
	var queries <-chan UrlQuery
	calls[0], queries, errors[0] = headPageQueriesDual[RespT](
		ctx,
		client,
		logger,
		initialQueries,
		callCap,
		queryCap,
		errorCap,
		totalPagesLimit,
	)
	calls[1], errors[1] = tailPageCallsDual[RespT](
		ctx,
		client,
		logger,
		queries,
		workersCap,
		errorCap,
	)
	return FanIn(calls), FanIn(errors)
}

func headPageQueriesDual[RespT any](
	ctx context.Context,
	client *Client,
	_ AppLog,
	firstQueries <-chan UrlQuery,
	callCap int,
	queryCap int,
	errorCap int,
	totalPagesLimit int,
) (
	<-chan CallNoError[RespT],
	<-chan UrlQuery,
	<-chan error,
) {
	calls := make(chan CallNoError[RespT], callCap)
	queries := make(chan UrlQuery, callCap)
	errors := make(chan error, errorCap)
	go func() {
		defer close(calls)
		defer close(queries)
		defer close(errors)
		for firstQuery := range firstQueries {
			firstVal, firstHeader, err := GetWithHeader[RespT](
				ctx,
				client,
				firstQuery.Path,
				firstQuery.Params)
			if err != nil {
				errors <- &UrlQueryError{err: err, query: firstQuery}
				continue
			}

			calls <- CallNoError[RespT]{
				Query:  firstQuery,
				Header: firstHeader,
				Val:    *firstVal,
			}

			cursor, err := parsePageCursor(firstHeader)
			if err != nil {
				errors <- &UrlQueryError{err: err, query: firstQuery}
				continue
			}
			if cursor.page == nil || cursor.totalPages == nil {
				continue
			}
			p := *cursor.page + 1
			n := *cursor.totalPages
			if totalPagesLimit > 0 && n > totalPagesLimit {
				n = totalPagesLimit
			}
			for ; p <= n; p++ {
				next := *firstQuery.Clone()
				next.Params["page"] = p
				queries <- next
			}
		}
	}()
	return calls, queries, errors
}

func tailPageCallsDual[RespT any](
	ctx context.Context,
	client *Client,
	_ AppLog,
	queries <-chan UrlQuery,
	workersCap int,
	errorsCap int,
) (<-chan CallNoError[RespT],
	<-chan error,
) {
	calls := make(chan CallNoError[RespT], workersCap)
	errors := make(chan error, errorsCap)
	go func() {
		defer close(calls)
		defer close(errors)
		var workersWg sync.WaitGroup
		workersWg.Add(workersCap)
		for i := 0; i < workersCap; i++ {
			go func() {
				defer workersWg.Done()
				for q := range queries {
					v, h, err := GetWithHeader[RespT](
						ctx,
						client,
						q.Path,
						q.Params)
					if err != nil {
						errors <- &UrlQueryError{err: err, query: q}
						continue
					}

					calls <- CallNoError[RespT]{
						Query:  q,
						Header: h,
						Val:    *v,
					}
				}
			}()
		}
		workersWg.Wait()

	}()
	return calls, errors
}
