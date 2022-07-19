package http

import (
	"net/http"

	"github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/logging"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func tracingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		otelhttp.NewHandler(next, "").ServeHTTP(w, r) // TODO what do we set operation name to?
	}
}

func logTracingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		traceID := span.SpanContext().TraceID()
		spanID := span.SpanContext().SpanID()

		ctx = logging.WithFields(ctx, zap.String("trace_id", traceID.String()), zap.String("span_id", spanID.String()))

		next(w, r.WithContext(ctx))
	}
}
