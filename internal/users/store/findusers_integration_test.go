package store_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
	"github.com/speakeasy-api/rest-template-go/internal/users/model"
	"github.com/speakeasy-api/rest-template-go/internal/users/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_FindUsers_Success(t *testing.T) {
	countries := []string{"UK", "IT", "US"}
	nicknames := []string{"nicky", "nicker", "nicksy", "nick"}
	firstNames := []string{"john", "judy", "jacob", "jacky", "jane"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Lee", "Clark"}
	domains := []string{"gmail.com", "hotmail.com", "yahoo.com"}

	for i := 0; i < 100; i++ {
		u := &model.User{
			FirstName: pointer.ToString(fmt.Sprintf("%s%d", firstNames[i%len(firstNames)], i)),
			LastName:  pointer.ToString(fmt.Sprintf("%s%d", lastNames[i%len(lastNames)], i)),
			Nickname:  pointer.ToString(fmt.Sprintf("%s%d", nicknames[i%len(nicknames)], i)),
			Password:  pointer.ToString("test"),
			Email:     pointer.ToString(fmt.Sprintf("%s%d@%s", firstNames[i%len(firstNames)], i, domains[i%len(domains)])),
			Country:   pointer.ToString(countries[i%len(countries)]),
		}

		_, err := insertUser(context.Background(), u)
		require.NoError(t, err)
	}

	type args struct {
		filters []model.Filter
		offset  int64
		limit   int64
	}
	tests := []struct {
		name          string
		args          args
		wantUserCount int
	}{
		{
			name: "get all UK users using equal match",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldCountry,
						MatchType: model.MatchTypeEqual,
						Value:     "UK",
					},
				},
				offset: 0,
				limit:  0,
			},
			wantUserCount: 35,
		},
		{
			name: "get all UK users using like match",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldCountry,
						MatchType: model.MatchTypeLike,
						Value:     "UK",
					},
				},
				offset: 0,
				limit:  0,
			},
			wantUserCount: 35,
		},
		{
			name: "get first page of all UK users using equal match",
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
			wantUserCount: 10,
		},
		{
			name: "get last page of all UK users using equal match",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldCountry,
						MatchType: model.MatchTypeEqual,
						Value:     "UK",
					},
				},
				offset: 30,
				limit:  10,
			},
			wantUserCount: 5,
		},
		{
			name: "get all users with a first name like john",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldFirstName,
						MatchType: model.MatchTypeLike,
						Value:     "john",
					},
				},
				offset: 0,
				limit:  0,
			},
			wantUserCount: 20,
		},
		{
			name: "get all users with a first name like john in the UK",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldFirstName,
						MatchType: model.MatchTypeLike,
						Value:     "john",
					},
					{
						Field:     model.FieldCountry,
						MatchType: model.MatchTypeEqual,
						Value:     "UK",
					},
				},
				offset: 0,
				limit:  0,
			},
			wantUserCount: 7,
		},
		{
			name: "get all users with a last name like Smith",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldLastName,
						MatchType: model.MatchTypeLike,
						Value:     "Smith",
					},
				},
				offset: 0,
				limit:  0,
			},
			wantUserCount: 17,
		},
		{
			name: "get all users with a last name like smith insensitive",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldLastName,
						MatchType: model.MatchTypeLike,
						Value:     "smith",
					},
				},
				offset: 0,
				limit:  0,
			},
			wantUserCount: 17,
		},
		{
			name: "get all users with a nickname like nicksy",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldNickname,
						MatchType: model.MatchTypeLike,
						Value:     "nicksy",
					},
				},
				offset: 0,
				limit:  0,
			},
			wantUserCount: 25,
		},
		{
			name: "get all users with an email like hotmail.com",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldEmail,
						MatchType: model.MatchTypeLike,
						Value:     "hotmail.com",
					},
				},
				offset: 0,
				limit:  0,
			},
			wantUserCount: 33,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			users, err := s.FindUsers(ctx, tt.args.filters, tt.args.offset, tt.args.limit)
			assert.NoError(t, err)
			assert.NotNil(t, users)
			assert.Len(t, users, tt.wantUserCount)
		})
	}
}

func TestStore_FindUsers_Error(t *testing.T) {
	type args struct {
		filters []model.Filter
		offset  int64
		limit   int64
	}
	tests := []struct {
		name     string
		args     args
		wantErr1 error
		wantErr2 error
	}{
		{
			name: "fails with no filters",
			args: args{
				filters: []model.Filter{},
				offset:  0,
				limit:   0,
			},
			wantErr1: errors.ErrInvalidRequest,
			wantErr2: store.ErrInvalidFilters,
		},
		{
			name: "fails with no users found",
			args: args{
				filters: []model.Filter{
					{
						Field:     model.FieldNickname,
						MatchType: model.MatchTypeLike,
						Value:     "blah",
					},
				},
				offset: 0,
				limit:  0,
			},
			wantErr1: errors.ErrNotFound,
			wantErr2: errors.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			users, err := s.FindUsers(ctx, tt.args.filters, tt.args.offset, tt.args.limit)
			assert.ErrorIs(t, err, tt.wantErr1)
			assert.ErrorIs(t, err, tt.wantErr2)
			assert.Nil(t, users)
		})
	}
}
