package fixtures

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

type MockSpan struct {
	// otel requires embed this interface to avoid a compilation errors on a private method span()
	// Apparently this enables avoiding breaking changes in the interface without package version changes
	embedded.Span // weird requirement
}

var _ trace.Span = (*MockSpan)(nil)

func (m *MockSpan) End(options ...trace.SpanEndOption) {}

func (m *MockSpan) AddEvent(name string, options ...trace.EventOption) {}

func (m *MockSpan) AddLink(link trace.Link) {}

func (m *MockSpan) IsRecording() bool { return false }

func (m *MockSpan) RecordError(err error, options ...trace.EventOption) {}

func (m *MockSpan) SpanContext() trace.SpanContext { return trace.SpanContext{} }

func (m *MockSpan) SetStatus(code codes.Code, description string) {}

func (m *MockSpan) SetName(name string) {}

func (m *MockSpan) SetAttributes(kv ...attribute.KeyValue) {}

func (m *MockSpan) TracerProvider() trace.TracerProvider { return nil }
