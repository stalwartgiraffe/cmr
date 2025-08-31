// Package fixtures provide testing helpers for package app
package fixtures

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

type MockApp struct {
	InitErr error
	MockLogger
}

func NewApp() *MockApp {
	return &MockApp{}
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

type MockLogger struct {
	SB strings.Builder
}

func (m *MockLogger) Printf(format string, v ...any) {
	m.SB.WriteString(fmt.Sprintf(format, v...))
}

func (m *MockLogger) Print(v ...any) {
	m.SB.WriteString(fmt.Sprint(v...))
}
func (m *MockLogger) Println(v ...any) {
	m.SB.WriteString(fmt.Sprintln(v...))
}
