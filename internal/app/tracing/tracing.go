package tracing

import (
	"context"
	"faceittechtest/internal/app/logging"
	"io"
	"io/ioutil"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tp *trace.TracerProvider

type OnShutdowner interface {
	OnShutdown(onShutdown func())
}

func EnableTracing(ctx context.Context, appName string, s OnShutdowner) error {
	exp, err := newExporter(ioutil.Discard) // TODO maybe be want this writen to a file or prometheus at some point
	if err != nil {
		return err
	}

	tp = trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource(appName)),
	)
	s.OnShutdown(func() {
		logging.From(ctx).Info("shutting down tracing provider")

		if err := tp.Shutdown(ctx); err != nil {
			logging.From(ctx).Error("failed to shutdown tracing provider")
		}
	})
	otel.SetTracerProvider(tp)

	return nil
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func newResource(appName string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
			//semconv.ServiceVersionKey.String("v0.1.0"), // TODO get version
			//attribute.String("environment", "demo"), // TODO determine other attributes
		),
	)
	return r
}
