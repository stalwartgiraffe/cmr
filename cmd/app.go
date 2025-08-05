package cmd

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type App interface {
	Tracer
}

type Tracer interface { 
	StartSpan(
		ctx context.Context,
		spanName string,
		opts ...trace.SpanStartOption) (
		context.Context,
		trace.Span)
}
