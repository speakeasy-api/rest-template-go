package tracing

import (
	"context"
	"io"

	"github.com/speakeasy-api/rest-template-go/internal/core/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tp *trace.TracerProvider

// OnShutdowner is an interface that allows a caller to register a function to be called when the application is shutting down.
type OnShutdowner interface {
	OnShutdown(onShutdown func())
}

// EnableTracing enables tracing.
func EnableTracing(ctx context.Context, appName string, s OnShutdowner) error {
	exp, err := stdouttrace.New(
		stdouttrace.WithWriter(io.Discard),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
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

func newResource(appName string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		),
	)
	return r
}
