package http

import "net/http"

type MiddlewareFunc func(h http.HandlerFunc) http.HandlerFunc

type handler struct {
	handler http.Handler
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware := []MiddlewareFunc{
		tracingMiddleware,
		logTracingMiddleware,
		requestLoggingMiddleware,
	}

	f := h.handler.ServeHTTP
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		f = m(f)
	}
	f(w, r)
}
