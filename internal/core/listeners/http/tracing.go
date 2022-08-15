package http

import (
	"net/http"

	"github.com/speakeasy-api/rest-template-go/internal/core/logging"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func tracingMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "example-rest-service")
}

func logTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		traceID := span.SpanContext().TraceID()
		spanID := span.SpanContext().SpanID()

		ctx = logging.WithFields(ctx, zap.String("trace_id", traceID.String()), zap.String("span_id", spanID.String()))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
