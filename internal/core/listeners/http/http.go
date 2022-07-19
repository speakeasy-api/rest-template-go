package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/errors"
	"github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/logging"

	"github.com/gorilla/mux"
)

const (
	ErrAddRoutes = errors.Error("failed to add routes")
	ErrServer    = errors.Error("listen stopped with error")
)

type Config struct {
	Port string `yaml:"port"`
}

type Service interface {
	AddRoutes(r *mux.Router) error
}

type Server struct {
	server *http.Server
	port   string
}

func New(s Service, cfg Config) (*Server, error) {
	r := mux.NewRouter()

	err := s.AddRoutes(r)
	if err != nil {
		return nil, ErrAddRoutes.Wrap(err)
	}

	h := handler{
		handler: r,
	}

	return &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", cfg.Port),
			BaseContext: func(net.Listener) context.Context {
				baseContext := context.Background()
				return logging.With(baseContext, logging.From(baseContext))
			},
			Handler: h,
		},
		port: cfg.Port,
	}, nil
}

func (s *Server) Listen(ctx context.Context) error {
	logging.From(ctx).Info(fmt.Sprintf("http server starting on port: %s", s.port))

	err := s.server.ListenAndServe()
	if err != nil {
		return ErrServer.Wrap(err)
	}

	logging.From(ctx).Info("http server stopped")

	return nil
}
