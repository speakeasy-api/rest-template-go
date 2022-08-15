package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
	"github.com/speakeasy-api/rest-template-go/internal/core/logging"
)

const (
	// ErrAddRoutes is the error returned when adding routes to the router fails.
	ErrAddRoutes = errors.Error("failed to add routes")
	// ErrServer is the error returned when the server stops due to an error.
	ErrServer = errors.Error("listen stopped with error")
)

const (
	readHeaderTimeout = 60 * time.Second
)

// Config represents the configuration of the http listener.
type Config struct {
	Port string `yaml:"port"`
}

// Service represents a http service that provides routes for the listener.
type Service interface {
	AddRoutes(r *mux.Router) error
}

// Server represents a http server that listens on a port.
type Server struct {
	server *http.Server
	port   string
}

// New instantiates a new instance of Server.
func New(s Service, cfg Config) (*Server, error) {
	r := mux.NewRouter()
	r.Use(tracingMiddleware)
	r.Use(logTracingMiddleware)
	r.Use(requestLoggingMiddleware)

	if err := s.AddRoutes(r); err != nil {
		return nil, ErrAddRoutes.Wrap(err)
	}

	return &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", cfg.Port),
			BaseContext: func(net.Listener) context.Context {
				baseContext := context.Background()
				return logging.With(baseContext, logging.From(baseContext))
			},
			Handler:           r,
			ReadHeaderTimeout: readHeaderTimeout,
		},
		port: cfg.Port,
	}, nil
}

// Listen starts the server and listens on the configured port.
func (s *Server) Listen(ctx context.Context) error {
	logging.From(ctx).Info(fmt.Sprintf("http server starting on port: %s", s.port))

	err := s.server.ListenAndServe()
	if err != nil {
		return ErrServer.Wrap(err)
	}

	logging.From(ctx).Info("http server stopped")

	return nil
}
