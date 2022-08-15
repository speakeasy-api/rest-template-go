package http

import (
	"net/http"

	"github.com/speakeasy-api/rest-template-go/internal/core/logging"
	"go.uber.org/zap"
)

func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithFields(r.Context(), zap.String("uri", r.RequestURI))

		logging.From(ctx).Info("request", zap.String("method", r.Method))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
