package store

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
	"github.com/speakeasy-api/rest-template-go/internal/users/model"
)

const (
	// ErrInvalidEmail is returned when the email is not a valid address or is empty.
	ErrInvalidEmail = errors.Error("invalid_email: email is invalid")
	// ErrEmailAlreadyUsed is returned when the email address is already used via another user.
	ErrEmailAlreadyUsed = errors.Error("email_already_used: email is already in use")
	// ErrEmptyNickname is returned when the nickname is empty.
	ErrEmptyNickname = errors.Error("empty_nickname: nickname is empty")
	// ErrNicknameAlreadyUsed is returned when the nickname is already used via another user.
	ErrNicknameAlreadyUsed = errors.Error("nickname_already_used: nickname is already in use")
	// ErrEmptyPassword is returned when the password is empty.
	ErrEmptyPassword = errors.Error("empty_password: password is empty")
	// ErrEmptyCountry is returned when the country is empty.
	ErrEmptyCountry = errors.Error("empty_country: password is empty")
	// ErrInvalidID si returned when the ID is not a valid UUID or is empty.
	ErrInvalidID = errors.Error("invalid_id: id is invalid")
	// ErrUserNotUpdated is returned when a record can't be found to update.
	ErrUserNotUpdated = errors.Error("user_not_updated: user record wasn't updated")
	// ErrUserNotDeleted is returned when a record can't be found to delete.
	ErrUserNotDeleted = errors.Error("user_not_deleted: user record wasn't deleted")
	// ErrInvalidFilters is returned when the filters for finding a user are not valid.
	ErrInvalidFilters = errors.Error("invalid_filters: filters invalid for finding user")
)

const (
	pqErrInvalidTextRepresentation = "invalid_text_representation"
)

var timeNow = func() *time.Time {
	now := time.Now().UTC()
	return &now
}

// DB represents a type for interfacing with a postgres database.
type DB interface {
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

// Store provides functionality for working with a postgres database.
type Store struct {
	db DB
}

// New will instantiate a new instance of Store.
func New(db DB) *Store {
	return &Store{
		db: db,
	}
}

// InsertUser will add a new unique user to the database using the provided data.
func (s *Store) InsertUser(ctx context.Context, u *model.User) (*model.User, error) {
	u.CreatedAt = timeNow()
	u.UpdatedAt = u.CreatedAt

	res, err := s.db.NamedQueryContext(ctx,
		`INSERT INTO 
		users(first_name, last_name, nickname, password, email, country, created_at, updated_at) 
		VALUES (:first_name, :last_name, :nickname, :password, :email, :country, :created_at, :updated_at) 
		RETURNING *`, u)
	if err = checkWriteError(err); err != nil {
		return nil, err
	}
	defer res.Close()

	if !res.Next() {
		return nil, errors.ErrUnknown
	}

	createdUser := &model.User{}

	if err := res.StructScan(&createdUser); err != nil {
		return nil, errors.ErrUnknown.Wrap(err)
	}

	return createdUser, nil
}

// GetUser will retrieve an existing user via their ID.
func (s *Store) GetUser(ctx context.Context, id string) (*model.User, error) {
	var u model.User

	if err := s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id = $1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Wrap(err)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == pqErrInvalidTextRepresentation && strings.Contains(pqErr.Error(), "uuid") {
				return nil, ErrInvalidID.Wrap(errors.ErrValidation.Wrap(err))
			}
		}

		return nil, errors.ErrUnknown.Wrap(err)
	}

	return &u, nil
}

// GetUserByEmail will retrieve an existing user via their email address.
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User

	if err := s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email = $1", email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Wrap(err)
		}

		return nil, errors.ErrUnknown.Wrap(err)
	}

	return &u, nil
}

// UpdateUser will update an existing user in the database using only the present data provided.
func (s *Store) UpdateUser(ctx context.Context, u *model.User) (*model.User, error) {
	if u.ID == nil || *u.ID == "" {
		return nil, ErrInvalidID.Wrap(errors.ErrValidation)
	}

	u.UpdatedAt = timeNow()

	res, err := s.db.NamedQueryContext(ctx,
		`UPDATE users 
		SET 
		first_name = COALESCE(:first_name, first_name), 
		last_name = COALESCE(:last_name, last_name), 
		nickname = COALESCE(:nickname, nickname), 
		password = COALESCE(:password, password),
		email = COALESCE(:email, email),
		country = COALESCE(:country, country),
		updated_at = :updated_at 
		WHERE id = :id
		RETURNING *`, u)
	if err = checkWriteError(err); err != nil {
		return nil, err
	}
	defer res.Close()

	if !res.Next() {
		return nil, ErrUserNotUpdated.Wrap(errors.ErrNotFound)
	}

	updatedUser := &model.User{}

	if err := res.StructScan(&updatedUser); err != nil {
		return nil, errors.ErrUnknown.Wrap(err)
	}

	return updatedUser, nil
}

// DeleteUser will delete an existing user via their ID.
func (s *Store) DeleteUser(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == pqErrInvalidTextRepresentation && strings.Contains(pqErr.Error(), "uuid") {
				return ErrInvalidID.Wrap(errors.ErrValidation.Wrap(err))
			}
		}

		return errors.ErrUnknown.Wrap(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return errors.ErrUnknown.Wrap(err)
	}
	if rows != 1 {
		return ErrUserNotDeleted.Wrap(errors.ErrNotFound)
	}

	return nil
}

//nolint:cyclop
func checkWriteError(err error) error {
	if err == nil {
		return nil
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code.Name() {
		case "string_data_right_truncation":
			return errors.ErrValidation.Wrap(err)
		case "check_violation":
			switch {
			case strings.Contains(pqErr.Error(), "email_check"):
				return ErrInvalidEmail.Wrap(errors.ErrValidation.Wrap(err))
			case strings.Contains(pqErr.Error(), "users_nickname_check"):
				return ErrEmptyNickname.Wrap(errors.ErrValidation.Wrap(err))
			case strings.Contains(pqErr.Error(), "users_password_check"):
				return ErrEmptyPassword.Wrap(errors.ErrValidation.Wrap(err))
			case strings.Contains(pqErr.Error(), "users_country_check"):
				return ErrEmptyCountry.Wrap(errors.ErrValidation.Wrap(err))
			default:
				return errors.ErrValidation.Wrap(err)
			}
		case "not_null_violation":
			switch {
			case strings.Contains(pqErr.Error(), "email"):
				return ErrInvalidEmail.Wrap(errors.ErrValidation.Wrap(err))
			case strings.Contains(pqErr.Error(), "nickname"):
				return ErrEmptyNickname.Wrap(errors.ErrValidation.Wrap(err))
			case strings.Contains(pqErr.Error(), "password"):
				return ErrEmptyPassword.Wrap(errors.ErrValidation.Wrap(err))
			case strings.Contains(pqErr.Error(), "country"):
				return ErrEmptyCountry.Wrap(errors.ErrValidation.Wrap(err))
			default:
				return errors.ErrValidation.Wrap(err)
			}
		case "unique_violation":
			if strings.Contains(pqErr.Error(), "email_unique") {
				return ErrEmailAlreadyUsed.Wrap(errors.ErrValidation.Wrap(err))
			} else if strings.Contains(pqErr.Error(), "nickname_unique") {
				return ErrNicknameAlreadyUsed.Wrap(errors.ErrValidation.Wrap(err))
			}
			return errors.ErrValidation.Wrap(err)
		case "invalid_text_representation":
			if strings.Contains(pqErr.Error(), "uuid") {
				return ErrInvalidID.Wrap(errors.ErrValidation.Wrap(err))
			}
		}
	}

	return errors.ErrUnknown.Wrap(err)
}
