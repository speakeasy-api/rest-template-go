package http

import (
	"net/http"

	"github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/logging"

	"go.uber.org/zap"
)

func requestLoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithFields(r.Context(), zap.String("uri", r.RequestURI))

		logging.From(ctx).Info("request", zap.String("method", r.Method)) // TODO determine how we might control request logging

		next(w, r.WithContext(ctx))
	}
}
