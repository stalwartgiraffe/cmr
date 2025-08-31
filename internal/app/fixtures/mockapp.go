// Package fixtures provide testing helpers for package app
package fixtures

import (
	"context"


	"go.opentelemetry.io/otel/trace"
)

type MockApp struct {
	InitErr error
}

func (a *MockApp) Err() error {
	return a.InitErr
}

func (a *MockApp) WithOtel(ctx context.Context, schema string) *MockApp {
	return a
}

func (a *MockApp) StartSpan(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption) (
	context.Context,
	trace.Span) {

	return ctx, &MockSpan{}
}
