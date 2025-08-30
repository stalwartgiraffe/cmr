package otel

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	vendorotel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"

	//"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Otel struct {
	Tracer trace.Tracer
}

// StartOtelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func StartOtelSDK(
	ctx context.Context,
	schemaName string,
	shutdowns ctxShutdowns,
) (
	Otel,
	error) {
	initPropagation()

	o := Otel{}
	tracer, err := startTracer(ctx, schemaName, shutdowns)
	// to export logs with an slog interface
	// logger = otelslog.NewLogger(schemaName)
	if err != nil {
		return o, err
	}
	//if err := startStdoutMetrics(ctx, shutdowns); err != nil {
	//	return o, err
	//}

	/* do we want otel logging export?
	if err := startLogger(ctx, shutdowns); err != nil {
		return o, err
	}
	*/

	/*
		configure how often we poll the go run time memstats
		if err := startReadMem(ctx, shutdowns); err != nil {
			return o, err
		}
	*/

	o.Tracer = tracer
	return o, nil
}

type ctxShutdowns interface {
	AddShutdown(f func(context.Context) error)
}

func initPropagation() {
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	vendorotel.SetTextMapPropagator(prop)
}

func startTracer(ctx context.Context, schemaName string, shutdowns ctxShutdowns) (trace.Tracer, error) {
	// this is the default client which expect to POST to an https server
	//traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient())

	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"), // Explicit HTTP with the default port number
		otlptracehttp.WithInsecure(),                 // Disable TLS
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(traceExporter))
	shutdowns.AddShutdown(tracerProvider.Shutdown)
	vendorotel.SetTracerProvider(tracerProvider)
	tracer := vendorotel.Tracer(schemaName)
	return tracer, nil
}

// this get go runtime metrics and dumps them to stdout
// for example "go.memory.used"
func startStdoutMetrics(ctx context.Context, shutdowns ctxShutdowns) error {
	// Print with a JSON encoder that indents with two spaces.
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	metricsExporter, err := stdoutmetric.New(
		stdoutmetric.WithEncoder(enc),
		// stdoutmetric.WithoutTimestamps(),
	)
	if err != nil {
		return err
	}

	res := resource.NewSchemaless(
		semconv.ServiceName("stdoutmetric-example"),
	)
	// Register the exporter with an SDK via a periodic reader.
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricsExporter)),
	)

	//ctx := context.Background()
	// This is where the sdk would be used to create a Meter and from that
	// instruments that would make measurements of your code. To simulate that
	// behavior, call export directly with mocked data.
	//_ = metricsExporter.Export(ctx, &mockData)

	// Ensure the periodic reader is cleaned up by shutting down the sdk.
	//_ = meterProvider.Shutdown(ctx)

	shutdowns.AddShutdown(meterProvider.Shutdown)
	vendorotel.SetMeterProvider(meterProvider)
	return nil
}

func startHttpMetrics(ctx context.Context, shutdowns ctxShutdowns) error {
	metricExporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return err
	}
	meterProvider :=
		metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(metricExporter)))

	shutdowns.AddShutdown(meterProvider.Shutdown)
	vendorotel.SetMeterProvider(meterProvider)
	return nil
}

func startHttpLogger(ctx context.Context, shutdowns ctxShutdowns) error {
	logExporter, err := otlploghttp.New(ctx, otlploghttp.WithInsecure()) // claude-ignore: required workaround
	if err != nil {
		return err
	}

	loggerProvider := log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(logExporter)))
	shutdowns.AddShutdown(loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	return nil
}

func startReadMem(ctx context.Context, shutdowns ctxShutdowns) error {
	return runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
}
