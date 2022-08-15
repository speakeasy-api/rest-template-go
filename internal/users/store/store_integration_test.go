package store_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/ory/dockertest/v3"
	"github.com/speakeasy-api/rest-template-go/internal/core/drivers/psql"
	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
	"github.com/speakeasy-api/rest-template-go/internal/users/model"
	"github.com/speakeasy-api/rest-template-go/internal/users/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	dsn = "postgresql://guest:guest@localhost:%s/speakeasy?sslmode=disable"
)

var (
	db                    *psql.Driver
	initialInsertedUserID string
)

var initialUser *model.User = &model.User{
	FirstName: pointer.ToString("testFirst"),
	LastName:  pointer.ToString("testLast"),
	Nickname:  pointer.ToString("test1"),
	Password:  pointer.ToString("test"),
	Email:     pointer.ToString("test1@test.com"),
	Country:   pointer.ToString("UK"),
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "alpine", []string{"POSTGRES_USER=guest", "POSTGRES_PASSWORD=guest", "POSTGRES_DB=speakeasy"})
	if err != nil {
		log.Fatalf("could not start resource: %v", err)
	}

	purge := func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("could not purge resource: %v", err)
		}
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		db = psql.New(psql.Config{
			DSN: fmt.Sprintf(dsn, resource.GetPort("5432/tcp")),
		})
		if err := db.Connect(ctx); err != nil {
			log.Printf("could not connect to database: %v", err)
			return err
		}

		if err := db.GetDB().Ping(); err != nil {
			log.Printf("could not ping to database: %v", err)
			return err
		}

		return nil
	}); err != nil {
		purge()
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := db.MigratePostgres(ctx, "file://../../../migrations"); err != nil {
		purge()
		log.Fatalf("could not migrate database: %v", err)
	}

	initialInsertedUserID, err = insertUser(ctx, initialUser)
	if err != nil {
		purge()
		log.Fatalf("could not insert user: %v", err)
	}

	code := m.Run()

	if err := db.Close(ctx); err != nil {
		purge()
		log.Fatalf("could not close db resource: %v", err)
	}

	purge()

	os.Exit(code)
}

func TestStore_InsertUser_Success(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name     string
		args     args
		wantUser model.User
	}{
		{
			name: "success",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("insertTest"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("inserttest@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantUser: model.User{
				FirstName: pointer.ToString("testFirst"),
				LastName:  pointer.ToString("testLast"),
				Nickname:  pointer.ToString("insertTest"),
				Password:  pointer.ToString("test"),
				Email:     pointer.ToString("inserttest@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			store.ExportSetTimeNow(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC))

			createdUser, err := s.InsertUser(ctx, tt.args.user)
			assert.NoError(t, err)
			require.NotNil(t, createdUser)

			tt.wantUser.ID = createdUser.ID

			assert.EqualValues(t, tt.wantUser, *createdUser)
		})
	}
}

func TestStore_InsertUser_Error(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name     string
		args     args
		wantErr1 error
		wantErr2 error
	}{
		{
			name: "failed with invalid email",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test2"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test2@@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrInvalidEmail,
		},
		{
			name: "failed with not-unique email",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test2"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test1@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrEmailAlreadyUsed,
		},
		{
			name: "failed with not-unique nickname",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test1"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test2@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrNicknameAlreadyUsed,
		},
		{
			name: "failed with empty email",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test2"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString(""),
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrInvalidEmail,
		},
		{
			name: "failed with null email",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test2"),
					Password:  pointer.ToString("test"),
					Email:     nil,
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrInvalidEmail,
		},
		{
			name: "failed with empty nickname",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString(""),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test2@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrEmptyNickname,
		},
		{
			name: "failed with null nickname",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  nil,
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test2@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrEmptyNickname,
		},
		{
			name: "failed with empty password",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test2"),
					Password:  pointer.ToString(""),
					Email:     pointer.ToString("test2@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrEmptyPassword,
		},
		{
			name: "failed with null password",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test2"),
					Password:  nil,
					Email:     pointer.ToString("test2@test.com"),
					Country:   pointer.ToString("UK"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrEmptyPassword,
		},
		{
			name: "failed with empty country",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test2"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test2@test.com"),
					Country:   pointer.ToString(""),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrEmptyCountry,
		},
		{
			name: "failed with null country",
			args: args{
				user: &model.User{
					FirstName: pointer.ToString("testFirst"),
					LastName:  pointer.ToString("testLast"),
					Nickname:  pointer.ToString("test2"),
					Password:  pointer.ToString("test"),
					Email:     pointer.ToString("test2@test.com"),
					Country:   nil,
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrEmptyCountry,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			createdUser, err := s.InsertUser(ctx, tt.args.user)
			assert.ErrorIs(t, err, tt.wantErr1)
			if tt.wantErr2 != nil {
				assert.ErrorIs(t, err, tt.wantErr2)
			}
			assert.Nil(t, createdUser)
		})
	}
}

func TestStore_GetUser_Success(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantUser model.User
	}{
		{
			name: "success",
			args: args{
				id: initialInsertedUserID,
			},
			wantUser: model.User{
				ID:        pointer.ToString(initialInsertedUserID),
				FirstName: pointer.ToString("testFirst"),
				LastName:  pointer.ToString("testLast"),
				Nickname:  pointer.ToString("test1"),
				Password:  pointer.ToString("test"),
				Email:     pointer.ToString("test1@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			user, err := s.GetUser(ctx, tt.args.id)
			assert.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, tt.wantUser, *user)
		})
	}
}

func TestStore_GetUser_Error(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantErr1 error
		wantErr2 error
	}{
		{
			name: "fail to find user",
			args: args{
				id: "9cdd8ae2-15ab-40df-9c46-50f391e16f60", // If we ever get a conflict this test will fail
			},
			wantErr1: errors.ErrNotFound,
			wantErr2: errors.ErrNotFound,
		},
		{
			name: "failed to get user with invalid id",
			args: args{
				id: "some-invalid-id",
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			user, err := s.GetUser(ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr1)
			assert.ErrorIs(t, err, tt.wantErr2)
			assert.Nil(t, user)
		})
	}
}

func TestStore_GetUserByEmail_Success(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name     string
		args     args
		wantUser model.User
	}{
		{
			name: "success",
			args: args{
				email: "test1@test.com",
			},
			wantUser: model.User{
				ID:        pointer.ToString(initialInsertedUserID),
				FirstName: pointer.ToString("testFirst"),
				LastName:  pointer.ToString("testLast"),
				Nickname:  pointer.ToString("test1"),
				Password:  pointer.ToString("test"),
				Email:     pointer.ToString("test1@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			user, err := s.GetUserByEmail(ctx, tt.args.email)
			assert.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, tt.wantUser, *user)
		})
	}
}

func TestStore_GetUserByEmail_Error(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "fail to find user by email",
			args: args{
				email: "not@found.com",
			},
			wantErr: errors.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			user, err := s.GetUserByEmail(ctx, tt.args.email)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, user)
		})
	}
}

func TestStore_UpdateUser_Success(t *testing.T) {
	updateTestUserID, err := insertUser(context.Background(), &model.User{
		FirstName: pointer.ToString("testFirst"),
		LastName:  pointer.ToString("testLast"),
		Nickname:  pointer.ToString("updateTest"),
		Password:  pointer.ToString("test"),
		Email:     pointer.ToString("updatetest@test.com"),
		Country:   pointer.ToString("UK"),
	})
	require.NoError(t, err)

	type fields struct {
		updateDay int
	}
	type args struct {
		user *model.User
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantUser model.User
	}{
		{
			name: "success updating first name",
			fields: fields{
				updateDay: 1,
			},
			args: args{
				user: &model.User{
					ID:        &updateTestUserID,
					FirstName: pointer.ToString("firstNameUpdate"),
				},
			},
			wantUser: model.User{
				ID:        &updateTestUserID,
				FirstName: pointer.ToString("firstNameUpdate"),
				LastName:  pointer.ToString("testLast"),
				Nickname:  pointer.ToString("updateTest"),
				Password:  pointer.ToString("test"),
				Email:     pointer.ToString("updatetest@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "success updating last name",
			fields: fields{
				updateDay: 2,
			},
			args: args{
				user: &model.User{
					ID:       &updateTestUserID,
					LastName: pointer.ToString("lastNameUpdate"),
				},
			},
			wantUser: model.User{
				ID:        &updateTestUserID,
				FirstName: pointer.ToString("firstNameUpdate"),
				LastName:  pointer.ToString("lastNameUpdate"),
				Nickname:  pointer.ToString("updateTest"),
				Password:  pointer.ToString("test"),
				Email:     pointer.ToString("updatetest@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2021, time.January, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "success updating nickname",
			fields: fields{
				updateDay: 3,
			},
			args: args{
				user: &model.User{
					ID:       &updateTestUserID,
					Nickname: pointer.ToString("nicknameUpdate"),
				},
			},
			wantUser: model.User{
				ID:        &updateTestUserID,
				FirstName: pointer.ToString("firstNameUpdate"),
				LastName:  pointer.ToString("lastNameUpdate"),
				Nickname:  pointer.ToString("nicknameUpdate"),
				Password:  pointer.ToString("test"),
				Email:     pointer.ToString("updatetest@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2021, time.January, 3, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "success updating password",
			fields: fields{
				updateDay: 4,
			},
			args: args{
				user: &model.User{
					ID:       &updateTestUserID,
					Password: pointer.ToString("passwordUpdate"),
				},
			},
			wantUser: model.User{
				ID:        &updateTestUserID,
				FirstName: pointer.ToString("firstNameUpdate"),
				LastName:  pointer.ToString("lastNameUpdate"),
				Nickname:  pointer.ToString("nicknameUpdate"),
				Password:  pointer.ToString("passwordUpdate"),
				Email:     pointer.ToString("updatetest@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2021, time.January, 4, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "success updating email",
			fields: fields{
				updateDay: 5,
			},
			args: args{
				user: &model.User{
					ID:    &updateTestUserID,
					Email: pointer.ToString("emailupdate@test.com"),
				},
			},
			wantUser: model.User{
				ID:        &updateTestUserID,
				FirstName: pointer.ToString("firstNameUpdate"),
				LastName:  pointer.ToString("lastNameUpdate"),
				Nickname:  pointer.ToString("nicknameUpdate"),
				Password:  pointer.ToString("passwordUpdate"),
				Email:     pointer.ToString("emailupdate@test.com"),
				Country:   pointer.ToString("UK"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2021, time.January, 5, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "success updating country",
			fields: fields{
				updateDay: 6,
			},
			args: args{
				user: &model.User{
					ID:      &updateTestUserID,
					Country: pointer.ToString("IT"),
				},
			},
			wantUser: model.User{
				ID:        &updateTestUserID,
				FirstName: pointer.ToString("firstNameUpdate"),
				LastName:  pointer.ToString("lastNameUpdate"),
				Nickname:  pointer.ToString("nicknameUpdate"),
				Password:  pointer.ToString("passwordUpdate"),
				Email:     pointer.ToString("emailupdate@test.com"),
				Country:   pointer.ToString("IT"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2021, time.January, 6, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "success updating everything",
			fields: fields{
				updateDay: 6,
			},
			args: args{
				user: &model.User{
					ID:        &updateTestUserID,
					FirstName: pointer.ToString("firstNameUpdate2"),
					LastName:  pointer.ToString("lastNameUpdate2"),
					Nickname:  pointer.ToString("nicknameUpdate2"),
					Password:  pointer.ToString("passwordUpdate2"),
					Email:     pointer.ToString("emailupdate2@test.com"),
					Country:   pointer.ToString("US"),
				},
			},
			wantUser: model.User{
				ID:        &updateTestUserID,
				FirstName: pointer.ToString("firstNameUpdate2"),
				LastName:  pointer.ToString("lastNameUpdate2"),
				Nickname:  pointer.ToString("nicknameUpdate2"),
				Password:  pointer.ToString("passwordUpdate2"),
				Email:     pointer.ToString("emailupdate2@test.com"),
				Country:   pointer.ToString("US"),
				CreatedAt: pointer.ToTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: pointer.ToTime(time.Date(2021, time.January, 6, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			store.ExportSetTimeNow(time.Date(2021, time.January, tt.fields.updateDay, 0, 0, 0, 0, time.UTC))

			updatedUser, err := s.UpdateUser(ctx, tt.args.user)
			assert.NoError(t, err)
			require.NotNil(t, updatedUser)
			assert.Equal(t, tt.wantUser, *updatedUser)
		})
	}
}

func TestStore_UpdateUser_Error(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name     string
		args     args
		wantErr1 error
		wantErr2 error
	}{
		{
			name: "failed with missing id",
			args: args{
				user: &model.User{
					ID:       nil,
					Nickname: pointer.ToString("testNickname"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrInvalidID,
		},
		{
			name: "failed with empty id",
			args: args{
				user: &model.User{
					ID:       pointer.ToString(""),
					Nickname: pointer.ToString("testNickname"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrInvalidID,
		},
		{
			name: "failed with invalid id",
			args: args{
				user: &model.User{
					ID:    pointer.ToString("invalid-id"),
					Email: pointer.ToString("updatetest@test.com"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrInvalidID,
		},
		{
			name: "failed with empty nickname",
			args: args{
				user: &model.User{
					ID:       &initialInsertedUserID,
					Nickname: pointer.ToString(""),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrEmptyNickname,
		},
		{
			name: "failed with empty password",
			args: args{
				user: &model.User{
					ID:       &initialInsertedUserID,
					Password: pointer.ToString(""),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrEmptyPassword,
		},
		{
			name: "failed with empty email",
			args: args{
				user: &model.User{
					ID:    &initialInsertedUserID,
					Email: pointer.ToString(""),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrInvalidEmail,
		},
		{
			name: "failed with invalid email",
			args: args{
				user: &model.User{
					ID:    &initialInsertedUserID,
					Email: pointer.ToString("test@@test.com"),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrInvalidEmail,
		},
		{
			name: "failed with empty country",
			args: args{
				user: &model.User{
					ID:      &initialInsertedUserID,
					Country: pointer.ToString(""),
				},
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrEmptyCountry,
		},
		{
			name: "failed updating non-existent user",
			args: args{
				user: &model.User{
					ID:    pointer.ToString("15b76e47-40df-43e6-9d1e-02b72c473914"), // Will fail if we get any UUID conflicts on pre inserted users
					Email: pointer.ToString("updatetest@test.com"),
				},
			},
			wantErr1: errors.ErrNotFound,
			wantErr2: store.ErrUserNotUpdated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			store.ExportSetTimeNow(time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC))

			updatedUser, err := s.UpdateUser(ctx, tt.args.user)
			assert.ErrorIs(t, err, tt.wantErr1)
			assert.ErrorIs(t, err, tt.wantErr2)
			assert.Nil(t, updatedUser)
		})
	}
}

func TestStore_DeleteUser_Success(t *testing.T) {
	deleteTestUserID, err := insertUser(context.Background(), &model.User{
		FirstName: pointer.ToString("testFirst"),
		LastName:  pointer.ToString("testLast"),
		Nickname:  pointer.ToString("deleteTest"),
		Password:  pointer.ToString("test"),
		Email:     pointer.ToString("deletetest@test.com"),
		Country:   pointer.ToString("UK"),
	})
	require.NoError(t, err)

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
				id: deleteTestUserID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			err := s.DeleteUser(ctx, tt.args.id)
			assert.NoError(t, err)
		})
	}
}

func TestStore_DeleteUser_Error(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantErr1 error
		wantErr2 error
	}{
		{
			name: "failed to delete non-existent user",
			args: args{
				id: "4f54b006-e7d9-47bf-ad38-d56c75a032cf",
			},
			wantErr1: errors.ErrNotFound,
			wantErr2: store.ErrUserNotDeleted,
		},
		{
			name: "failed to delete user with invalid id",
			args: args{
				id: "some-invalid-id",
			},
			wantErr1: errors.ErrValidation,
			wantErr2: store.ErrInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(db.GetDB())

			ctx := context.Background()

			err := s.DeleteUser(ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr1)
			assert.ErrorIs(t, err, tt.wantErr2)
		})
	}
}

func insertUser(ctx context.Context, u *model.User) (string, error) {
	s := store.New(db.GetDB())

	store.ExportSetTimeNow(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC))

	createdUser, err := s.InsertUser(ctx, u)
	if err != nil {
		return "", err
	}

	return *createdUser.ID, nil
}
