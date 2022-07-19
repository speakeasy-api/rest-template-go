package model

import "time"

// User represents a person using our platform.
type User struct {
	ID        *string    `json:"id" db:"id"`
	FirstName *string    `json:"first_name" db:"first_name"`
	LastName  *string    `json:"last_name" db:"last_name"`
	Nickname  *string    `json:"nickname" db:"nickname"`
	Password  *string    `json:"password" db:"password"`
	Email     *string    `json:"email" db:"email"`
	Country   *string    `json:"country" db:"country"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Field is an enum providing valid fields for filtering.
type Field string

const (
	// FieldFirstName represents the first name field.
	FieldFirstName Field = "first_name"
	// FieldLastName represents the last name field.
	FieldLastName Field = "last_name"
	// FieldNickname represents the nickname field.
	FieldNickname Field = "nickname"
	// FieldEmail represents the email field.
	FieldEmail Field = "email"
	// FieldCountry represents the country field.
	FieldCountry Field = "country"
)

// MatchType is an enum providing valid matching mechanisms for filtering values.
type MatchType string

const (
	// MatchTypeLike represents a LIKE match.
	MatchTypeLike MatchType = "ILIKE"
	// MatchTypeEqual represents an exact match.
	MatchTypeEqual MatchType = "="
)

// Filter is a struct representing a filter for finding users.
type Filter struct {
	MatchType MatchType `json:"match_type"`
	Field     Field     `json:"field"`
	Value     string    `json:"value"`
}
