package users_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
	"github.com/speakeasy-api/rest-template-go/internal/events"
	"github.com/speakeasy-api/rest-template-go/internal/users"
	"github.com/speakeasy-api/rest-template-go/internal/users/mocks"
	"github.com/speakeasy-api/rest-template-go/internal/users/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockStore(ctrl)
	e := mocks.NewMockEvents(ctrl)

	u := users.New(s, e)
	assert.NotNil(t, u)
}

func TestUsers_CreateUser_Success(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name     string
		args     args
		wantUser *model.User
	}{
		{
			name: "success",
			args: args{
				user: &model.User{
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
				CreatedAt: pointer.ToTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			e := mocks.NewMockEvents(ctrl)

			u := users.New(s, e)
			require.NotNil(t, u)

			ctx := context.Background()

			s.EXPECT().InsertUser(gomock.Any(), tt.args.user).Return(tt.wantUser, nil).Times(1)
			e.EXPECT().Produce(gomock.Any(), events.TopicUsers, events.UserEvent{
				EventType: events.EventTypeUserCreated,
				ID:        *tt.wantUser.ID,
				User:      tt.wantUser,
			}).Times(1)

			user, err := u.CreateUser(ctx, tt.args.user)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantUser, user)
		})
	}
}

func TestUsers_CreateUser_Error(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "fails",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr: errors.New("test fail"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			e := mocks.NewMockEvents(ctrl)

			u := users.New(s, e)
			require.NotNil(t, u)

			ctx := context.Background()

			s.EXPECT().InsertUser(gomock.Any(), tt.args.user).Return(nil, tt.wantErr).Times(1)

			user, err := u.CreateUser(ctx, tt.args.user)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, user)
		})
	}
}

func TestUsers_GetUser_Success(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantUser *model.User
	}{
		{
			name: "success",
			args: args{},
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			e := mocks.NewMockEvents(ctrl)

			u := users.New(s, e)
			require.NotNil(t, u)

			ctx := context.Background()

			s.EXPECT().GetUser(gomock.Any(), tt.args.id).Return(tt.wantUser, nil).Times(1)

			user, err := u.GetUser(ctx, tt.args.id)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.wantUser, user)
		})
	}
}

func TestUsers_GetUser_Error(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "fails",
			args: args{
				id: "some-test-id",
			},
			wantErr: errors.New("test fail"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			e := mocks.NewMockEvents(ctrl)

			u := users.New(s, e)
			require.NotNil(t, u)

			ctx := context.Background()

			s.EXPECT().GetUser(gomock.Any(), tt.args.id).Return(nil, tt.wantErr).Times(1)

			user, err := u.GetUser(ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, user)
		})
	}
}

func TestUsers_FindUsers_Success(t *testing.T) {
	type args struct {
		filters []model.Filter
		offset  int64
		limit   int64
	}
	tests := []struct {
		name      string
		args      args
		wantUsers []*model.User
	}{
		{
			name: "success",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldCountry,
						MatchType: model.MatchTypeEqual,
						Value:     "UK",
					},
				},
				offset: 0,
				limit:  10,
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			e := mocks.NewMockEvents(ctrl)

			u := users.New(s, e)
			require.NotNil(t, u)

			ctx := context.Background()

			s.EXPECT().FindUsers(gomock.Any(), tt.args.filters, tt.args.offset, tt.args.limit).Return(tt.wantUsers, nil).Times(1)

			users, err := u.FindUsers(ctx, tt.args.filters, tt.args.offset, tt.args.limit)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.wantUsers, users)
		})
	}
}

func TestUsers_FindUsers_Error(t *testing.T) {
	type fields struct {
		findUsersErr error
	}
	type args struct {
		filters []model.Filter
		offset  int64
		limit   int64
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr1 error
		wantErr2 error
	}{
		{
			name: "fails when search fails",
			fields: fields{
				findUsersErr: errors.ErrUnknown,
			},
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldCountry,
						MatchType: model.MatchTypeEqual,
						Value:     "UK",
					},
				},
				offset: 0,
				limit:  10,
			},
			wantErr1: errors.ErrUnknown,
			wantErr2: errors.ErrUnknown,
		},
		{
			name: "fails with empty value",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldCountry,
						MatchType: model.MatchTypeEqual,
						Value:     "",
					},
				},
				offset: 0,
				limit:  10,
			},
			wantErr1: errors.ErrValidation,
			wantErr2: users.ErrInvalidFilterValue,
		},
		{
			name: "fails with invalid match type",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldCountry,
						MatchType: "invalid",
						Value:     "UK",
					},
				},
				offset: 0,
				limit:  10,
			},
			wantErr1: errors.ErrValidation,
			wantErr2: users.ErrInvalidFilterMatchType,
		},
		{
			name: "fails with invalid field",
			args: args{
				filters: []model.Filter{
					{
						Field:     "invalid",
						MatchType: model.MatchTypeEqual,
						Value:     "UK",
					},
				},
				offset: 0,
				limit:  10,
			},
			wantErr1: errors.ErrValidation,
			wantErr2: users.ErrInvalidFilterField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			e := mocks.NewMockEvents(ctrl)

			u := users.New(s, e)
			require.NotNil(t, u)

			ctx := context.Background()

			if tt.fields.findUsersErr != nil {
				s.EXPECT().FindUsers(gomock.Any(), tt.args.filters, tt.args.offset, tt.args.limit).Return(nil, tt.fields.findUsersErr).Times(1)
			}

			user, err := u.FindUsers(ctx, tt.args.filters, tt.args.offset, tt.args.limit)
			assert.ErrorIs(t, err, tt.wantErr1)
			assert.ErrorIs(t, err, tt.wantErr2)
			assert.Nil(t, user)
		})
	}
}

func TestUsers_UpdateUser_Success(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name     string
		args     args
		wantUser *model.User
	}{
		{
			name: "success",
			args: args{
				user: &model.User{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			e := mocks.NewMockEvents(ctrl)

			u := users.New(s, e)
			require.NotNil(t, u)

			ctx := context.Background()

			s.EXPECT().UpdateUser(gomock.Any(), tt.args.user).Return(tt.wantUser, nil).Times(1)
			e.EXPECT().Produce(gomock.Any(), events.TopicUsers, events.UserEvent{
				EventType: events.EventTypeUserUpdated,
				ID:        *tt.wantUser.ID,
				User:      tt.wantUser,
			}).Times(1)

			updatedUser, err := u.UpdateUser(ctx, tt.args.user)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantUser, updatedUser)
		})
	}
}

func TestUsers_UpdateUser_Error(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "fails",
			args: args{
				user: &model.User{
					ID:        pointer.ToString("some-test-id"),
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr: errors.New("test fail"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			e := mocks.NewMockEvents(ctrl)

			u := users.New(s, e)
			require.NotNil(t, u)

			ctx := context.Background()

			s.EXPECT().UpdateUser(gomock.Any(), tt.args.user).Return(nil, tt.wantErr).Times(1)

			updatedUser, err := u.UpdateUser(ctx, tt.args.user)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, updatedUser)
		})
	}
}

func TestUsers_DeleteUser_Success(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				id: "some-test-id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			e := mocks.NewMockEvents(ctrl)

			u := users.New(s, e)
			require.NotNil(t, u)

			ctx := context.Background()

			s.EXPECT().DeleteUser(gomock.Any(), tt.args.id).Return(nil).Times(1)
			e.EXPECT().Produce(gomock.Any(), events.TopicUsers, events.UserEvent{
				EventType: events.EventTypeUserDeleted,
				ID:        tt.args.id,
			}).Times(1)

			err := u.DeleteUser(ctx, tt.args.id)
			assert.NoError(t, err)
		})
	}
}

func TestUsers_DeleteUser_Error(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "fails",
			args: args{
				id: "some-test-id",
			},
			wantErr: errors.New("test fail"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockStore(ctrl)
			e := mocks.NewMockEvents(ctrl)

			u := users.New(s, e)
			require.NotNil(t, u)

			ctx := context.Background()

			s.EXPECT().DeleteUser(gomock.Any(), tt.args.id).Return(tt.wantErr).Times(1)

			err := u.DeleteUser(ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
