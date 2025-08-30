package cmd

import (
	"context"

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
	Println(v ...any)
	Printf(format string, v ...any)
	Print(v ...any)
}
