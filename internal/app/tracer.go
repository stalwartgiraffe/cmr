// Package app provides injectable application level embedded singletons that can be passed as a function argument.
// To avoid circular dependencies, App is a concrete pointer. Consuming packages should define accessor interfaces.
package app

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/stalwartgiraffe/cmr/internal/otel"
)

func (c AppErr) WithOtel(ctx context.Context, schema string) AppErr {
	if c.Err != nil {
		return c
	}

	c.App.Otel, c.Err = otel.StartOtelSDK(ctx, schema, c.App)

	return c

}

func (a *App) StartSpan(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption) (
	context.Context,
	trace.Span) {
	return a.Tracer.Start(ctx, spanName, opts...)
}
