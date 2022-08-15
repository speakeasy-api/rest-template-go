package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	httptransport "github.com/speakeasy-api/rest-template-go/internal/transport/http"
	"github.com/speakeasy-api/rest-template-go/internal/transport/http/mocks"
	"github.com/speakeasy-api/rest-template-go/internal/users/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseUserURL = "/v1/user"
	userURL     = baseUserURL + "/%s"
	searchURL   = "/v1/users/search"
)

func TestServer_CreateUser_Success(t *testing.T) {
	type args struct {
		user model.User
	}
	tests := []struct {
		name     string
		args     args
		wantUser model.User
		wantCode int
	}{
		{
			name: "success",
			args: args{
				user: model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantUser: model.User{
				ID:        pointer.ToString("some-test-id"),
				FirstName: pointer.ToString("testFirst"),
				LastName:  pointer.ToString("testLast"),
				Nickname:  pointer.ToString("test"),
				Password:  pointer.ToString("test"),
				Email:     pointer.ToString("test@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
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

			u.EXPECT().CreateUser(gomock.Any(), &tt.args.user).Return(&tt.wantUser, nil).Times(1)

			data, err := json.Marshal(tt.args.user)
			require.NoError(t, err)
			require.NotNil(t, data)

			req, err := http.NewRequest(http.MethodPost, baseUserURL, bytes.NewBuffer(data))
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Data model.User `json:"data"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tt.wantUser, res.Data)
		})
	}
}

func TestServer_CreateUser_Error(t *testing.T) {
	type args struct {
		user model.User
	}
	tests := []struct {
		name     string
		args     args
		wantErr  string
		wantCode int
	}{
		{
			name: "fails",
			args: args{
				user: model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
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

			u.EXPECT().CreateUser(gomock.Any(), &tt.args.user).Return(nil, errors.New(tt.wantErr)).Times(1)

			data, err := json.Marshal(tt.args.user)
			require.NoError(t, err)
			require.NotNil(t, data)

			req, err := http.NewRequest(http.MethodPost, baseUserURL, bytes.NewBuffer(data))
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

func TestServer_GetUser_Success(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantUser model.User
		wantCode int
	}{
		{
			name: "success",
			args: args{
				id: "some-test-id",
			},
			wantUser: model.User{
				ID:        pointer.ToString("some-test-id"),
				FirstName: pointer.ToString("testFirst"),
				LastName:  pointer.ToString("testLast"),
				Nickname:  pointer.ToString("test"),
				Password:  pointer.ToString("test"),
				Email:     pointer.ToString("test@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
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

			u.EXPECT().GetUser(gomock.Any(), tt.args.id).Return(&tt.wantUser, nil).Times(1)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(userURL, tt.args.id), nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Data model.User `json:"data"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.EqualValues(t, tt.wantUser, res.Data)
		})
	}
}

func TestServer_GetUser_Error(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  string
		wantCode int
	}{
		{
			name: "fails",
			args: args{
				id: "some-test-id",
			},
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

			u.EXPECT().GetUser(gomock.Any(), tt.args.id).Return(nil, errors.New(tt.wantErr)).Times(1)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(userURL, tt.args.id), nil)
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

func TestServer_SearchUsers_Success(t *testing.T) {
	type args struct {
		filters []model.Filter
		offset  int64
		limit   int64
	}
	tests := []struct {
		name      string
		args      args
		wantUsers []*model.User
		wantCode  int
	}{
		{
			name: "success",
			args: args{
				filters: []model.Filter{
					{
						MatchType: model.MatchTypeEqual,
						Field:     model.FieldCountry,
						Value:     "UK",
					},
				},
			},
			wantUsers: []*model.User{
				{
					ID:        pointer.ToString("some-test-id"),
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test@test.com"),
					Country:   pointer.ToString("UK"),
					CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
					UpdatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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

			u.EXPECT().FindUsers(gomock.Any(), tt.args.filters, tt.args.offset, tt.args.limit).Return(tt.wantUsers, nil).Times(1)

			data, err := json.Marshal(httptransport.SearchUsersRequest{Filters: tt.args.filters, Offset: tt.args.offset, Limit: tt.args.limit})
			require.NoError(t, err)
			require.NotNil(t, data)

			req, err := http.NewRequest(http.MethodPost, searchURL, bytes.NewBuffer(data))
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Data []*model.User `json:"data"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.EqualValues(t, tt.wantUsers, res.Data)
		})
	}
}

func TestServer_SearchUsers_Error(t *testing.T) {
	type args struct {
		filters []model.Filter
		offset  int64
		limit   int64
	}
	tests := []struct {
		name     string
		args     args
		wantErr  string
		wantCode int
	}{
		{
			name: "success",
			args: args{
				filters: []model.Filter{
					{
						MatchType: model.MatchTypeEqual,
						Field:     model.FieldCountry,
						Value:     "UK",
					},
				},
			},
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

			u.EXPECT().FindUsers(gomock.Any(), tt.args.filters, tt.args.offset, tt.args.limit).Return(nil, errors.New(tt.wantErr)).Times(1)

			data, err := json.Marshal(httptransport.SearchUsersRequest{Filters: tt.args.filters, Offset: tt.args.offset, Limit: tt.args.limit})
			require.NoError(t, err)
			require.NotNil(t, data)

			req, err := http.NewRequest(http.MethodPost, searchURL, bytes.NewBuffer(data))
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

func TestServer_UpdateUser_Success(t *testing.T) {
	type args struct {
		user model.User
	}
	tests := []struct {
		name     string
		args     args
		wantUser *model.User
		wantCode int
	}{
		{
			name: "success",
			args: args{
				user: model.User{
					ID:        pointer.ToString("some-test-id"),
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantUser: &model.User{
				ID:        pointer.ToString("some-test-id"),
				FirstName: pointer.ToString("testFirst"),
				LastName:  pointer.ToString("testLast"),
				Nickname:  pointer.ToString("test"),
				Password:  pointer.ToString("test"),
				Email:     pointer.ToString("test@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
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

			u.EXPECT().UpdateUser(gomock.Any(), &tt.args.user).Return(tt.wantUser, nil).Times(1)

			data, err := json.Marshal(tt.args.user)
			require.NoError(t, err)
			require.NotNil(t, data)

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(userURL, *tt.args.user.ID), bytes.NewBuffer(data))
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Data *model.User `json:"data"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tt.wantUser, res.Data)
		})
	}
}

func TestServer_UpdateUser_Error(t *testing.T) {
	type args struct {
		user model.User
	}
	tests := []struct {
		name     string
		args     args
		wantErr  string
		wantCode int
	}{
		{
			name: "fails",
			args: args{
				user: model.User{
					ID:        pointer.ToString("some-test-id"),
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
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

			u.EXPECT().UpdateUser(gomock.Any(), &tt.args.user).Return(nil, errors.New(tt.wantErr)).Times(1)

			data, err := json.Marshal(tt.args.user)
			require.NoError(t, err)
			require.NotNil(t, data)

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(userURL, *tt.args.user.ID), bytes.NewBuffer(data))
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

func TestServer_DeleteUser_Success(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		wantSuccess bool
		wantCode    int
	}{
		{
			name: "success",
			args: args{
				id: "some-test-id",
			},
			wantSuccess: true,
			wantCode:    http.StatusOK,
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

			u.EXPECT().DeleteUser(gomock.Any(), tt.args.id).Return(nil).Times(1)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(userURL, tt.args.id), nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res struct {
				Data httptransport.DeletedUserResponse `json:"data"`
			}

			err = json.Unmarshal(w.Body.Bytes(), &res)
			require.NoError(t, err)
			assert.Equal(t, tt.wantSuccess, res.Data.Success)
		})
	}
}

func TestServer_DeleteUser_Error(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  string
		wantCode int
	}{
		{
			name: "fails",
			args: args{
				id: "some-test-id",
			},
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

			u.EXPECT().DeleteUser(gomock.Any(), tt.args.id).Return(errors.New(tt.wantErr)).Times(1)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(userURL, tt.args.id), nil)
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
