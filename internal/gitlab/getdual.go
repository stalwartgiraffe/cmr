package gitlab

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

type App interface {
	Tracer
	Logger
}

type Tracer interface {
	StartSpan(
		ctx context.Context,
		spanName string,
		opts ...trace.SpanStartOption) (
		context.Context,
		trace.Span)
}

type Logger interface {
	Printf(format string, v ...any)
	Print(v ...any)
	Println(v ...any)
}

func GatherPageCallsDualApp[RespT any](
	ctx context.Context,
	app App,
	client *Client,
	initialQueries <-chan UrlQuery,
	totalPagesLimit int,
) (
	<-chan CallNoError[RespT],
	<-chan error,
) {
	return GatherPageCallsWithDualApp[RespT](
		ctx,
		app,
		client,
		initialQueries,
		5,               // callCap int,
		5,               // queryCap int,
		5,               // workersCap int,
		1,               // errorCap int,
		totalPagesLimit, // 0 means no limit
	)
}

func GatherPageCallsDual[RespT any](
	ctx context.Context,
	app App,
	client *Client,
	initialQueries <-chan UrlQuery,
	totalPagesLimit int,
) (
	<-chan CallNoError[RespT],
	<-chan error,
) {
	return GatherPageCallsWithDual[RespT](
		ctx,
		app,
		client,
		initialQueries,
		5,               // callCap int,
		5,               // queryCap int,
		5,               // workersCap int,
		1,               // errorCap int,
		totalPagesLimit, // 0 means no limit
	)
}

func GatherPageCallsWithDualApp[RespT any](
	ctx context.Context,
	app App,
	client *Client,
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
	if app != nil {
		var span trace.Span
		ctx, span = app.StartSpan(ctx, "GatherPageCallsWithDualApp")
		defer span.End()
	}

	calls := make([]<-chan CallNoError[RespT], 2)
	errors := make([]<-chan error, 2)
	var queries <-chan UrlQuery
	calls[0], queries, errors[0] = headPageQueriesDual[RespT](
		ctx,
		app,
		client,
		initialQueries,
		callCap,
		queryCap,
		errorCap,
		totalPagesLimit,
	)
	calls[1], errors[1] = tailPageCallsDual[RespT](
		ctx,
		app,
		client,
		queries,
		workersCap,
		errorCap,
	)
	return FanIn(calls), FanIn(errors)
}

func GatherPageCallsWithDual[RespT any](
	ctx context.Context,
	app App,
	client *Client,
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
	return GatherPageCallsWithDualApp[RespT](
		ctx,
		app,
		client,
		initialQueries,
		callCap,
		queryCap,
		workersCap,
		errorCap,
		totalPagesLimit,
	)
}

func headPageQueriesDual[RespT any](
	ctx context.Context,
	app App,
	client *Client,
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
	ctx, span := app.StartSpan(ctx, "headPageQueriesDual")
	defer span.End()

	calls := make(chan CallNoError[RespT], callCap)
	queries := make(chan UrlQuery, callCap)
	errors := make(chan error, errorCap)
	go func() {
		ctx, span = app.StartSpan(ctx, "go_headPageQueriesDual")
		defer span.End()
		defer close(calls)
		defer close(queries)
		defer close(errors)
		for firstQuery := range firstQueries {
			firstVal, firstHeader, err := GetWithHeader[RespT](
				ctx,
				app,
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
	app App,
	client *Client,
	queries <-chan UrlQuery,
	workersCap int,
	errorsCap int,
) (<-chan CallNoError[RespT],
	<-chan error,
) {
	ctx, span := app.StartSpan(ctx, "tailPageCallsDuel")
	defer span.End()
	calls := make(chan CallNoError[RespT], workersCap)
	errors := make(chan error, errorsCap)
	go func() {
		ctx, span := app.StartSpan(ctx, "go_tailPageCallsDuel")
		defer span.End()
		defer close(calls)
		defer close(errors)
		var workersWg sync.WaitGroup
		workersWg.Add(workersCap)
		for i := 0; i < workersCap; i++ {
			go func() {
				ctx, span := app.StartSpan(ctx, "go_xo_tailPageCallsDuel")
				defer span.End()

				defer workersWg.Done()
				for q := range queries {
					v, h, err := GetWithHeader[RespT](
						ctx,
						app,
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
