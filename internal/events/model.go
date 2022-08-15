package events

import "github.com/speakeasy-api/rest-template-go/internal/users/model"

// EventType represents the type of event that occurred.
type EventType string

const (
	// EventTypeUserCreated is triggered after a user has been successfully created.
	EventTypeUserCreated EventType = "user_created"
	// EventTypeUserUpdated is triggered after a user has been successfully updated.
	EventTypeUserUpdated EventType = "user_updated"
	// EventTypeUserDeleted is triggered after a user has been successfully deleted.
	EventTypeUserDeleted EventType = "user_deleted"
)

// UserEvent represents an event that occurs on a user entity.
type UserEvent struct {
	EventType EventType   `json:"event_type"`
	ID        string      `json:"id"`
	User      *model.User `json:"user"`
}
