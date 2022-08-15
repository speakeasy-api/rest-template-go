package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	httptransport "github.com/speakeasy-api/rest-template-go/internal/transport/http"
	"github.com/speakeasy-api/rest-template-go/internal/transport/http/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_Health_Success(t *testing.T) {
	tests := []struct {
		name     string
		wantCode int
	}{
		{
			name:     "success",
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mocks.NewMockUsers(ctrl)
			d := mocks.NewMockDB(ctrl)

			ht := httptransport.New(u, d)
			require.NotNil(t, ht)

			r := mux.NewRouter()

			err := ht.AddRoutes(r)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			d.EXPECT().PingContext(gomock.Any()).Return(nil).Times(1)

			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestServer_Health_Error(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  string
		wantCode int
	}{
		{
			name:     "fails",
			wantErr:  "test fail",
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mocks.NewMockUsers(ctrl)
			d := mocks.NewMockDB(ctrl)

			ht := httptransport.New(u, d)
			require.NotNil(t, ht)

			r := mux.NewRouter()

			err := ht.AddRoutes(r)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			d.EXPECT().PingContext(gomock.Any()).Return(errors.New(tt.wantErr)).Times(1)

			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Error string `json:"error"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tt.wantErr, res.Error)
		})
	}
}
