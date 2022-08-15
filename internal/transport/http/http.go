//go:generate mockgen -destination=./mocks/http_mock.go -package mocks github.com/speakeasy-api/rest-template-go/internal/transport/http Users,DB

package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/speakeasy-api/rest-template-go/internal/users/model"
)

// Users represents a type that can provide CRUD operations on users.
type Users interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	FindUsers(ctx context.Context, filters []model.Filter, offset, limit int64) ([]*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

// DB represents a type that can be used to interact with the database.
type DB interface {
	PingContext(ctx context.Context) error
}

// Server represents a HTTP server that can handle requests for this microservice.
type Server struct {
	users Users
	db    DB
}

// New will instantiate a new instance of Server.
func New(u Users, db DB) *Server {
	return &Server{
		users: u,
		db:    db,
	}
}

// AddRoutes will add the routes this server supports to the router.
func (s *Server) AddRoutes(r *mux.Router) error {
	r.HandleFunc("/health", s.healthCheck).Methods(http.MethodGet)

	r = r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/user", s.createUser).Methods(http.MethodPost)
	r.HandleFunc("/user/{id}", s.getUser).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", s.updateUser).Methods(http.MethodPut)
	r.HandleFunc("/user/{id}", s.deleteUser).Methods(http.MethodDelete)

	// Not the most RESTful way of doing this as it won't really be cachable but provides easier parsing of the inputs for now
	r.HandleFunc("/users/search", s.searchUsers).Methods(http.MethodPost)

	return nil
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.db.PingContext(r.Context()); err != nil {
		handleError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	jsonRes := struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}

	dataBytes, err := json.Marshal(jsonRes)
	if err != nil {
		handleError(ctx, w, err)
		return
	}

	if _, err := w.Write(dataBytes); err != nil {
		handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
