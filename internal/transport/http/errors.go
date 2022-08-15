package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
	"github.com/speakeasy-api/rest-template-go/internal/core/logging"
	"go.uber.org/zap"
)

func handleError(ctx context.Context, w http.ResponseWriter, err error) {
	logging.From(ctx).Error("error occurred in request", zap.Error(err))

	switch {
	case errors.Is(err, errors.ErrInvalidRequest):
		fallthrough
	case errors.Is(err, errors.ErrValidation):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, errors.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, errors.ErrUnknown):
		fallthrough
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	errJSON := struct {
		Error string `json:"error"`
	}{
		Error: strings.Split(err.Error(), errors.ErrSeperator)[0], // TODO we may need to strip additional error information
	}

	data, err := json.Marshal(errJSON)
	if err != nil {
		logging.From(ctx).Error("failed to serialize error response", zap.Error(err))
		data = []byte(`{"error": "internal server error"}`)
	}

	_, err = w.Write(data)
	if err != nil {
		logging.From(ctx).Error("failed to write error response", zap.Error(err))
	}
}
