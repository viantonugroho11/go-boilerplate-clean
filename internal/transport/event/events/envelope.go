package events

import (
	"time"

	"github.com/google/uuid"
)

type Action string

const (
	ActionCreate Action = "CREATE"
	ActionUpdate Action = "UPDATE"
	ActionDelete Action = "DELETE"
)

type Meta struct {
	EventID              string `json:"event_id"`
	EventTimestamp       string `json:"event_timestamp"`
	Action               Action `json:"action"`
	Resource             string `json:"resource"`
	MessageSchemaVersion int    `json:"message_schema_version"`
}

// Event is the standard CDC envelope for all domain events.
type Event[T any] struct {
	ResourceID string `json:"resource_id"`
	Meta       Meta   `json:"meta"`
	Before     *T     `json:"before,omitempty"`
	After      *T     `json:"after,omitempty"`
}

func NewEvent[T any](resourceID, resource string, action Action, before, after *T) Event[T] {
	return Event[T]{
		ResourceID: resourceID,
		Meta: Meta{
			EventID:              uuid.NewString(),
			EventTimestamp:       time.Now().UTC().Format(time.RFC3339),
			Action:               action,
			Resource:             resource,
			MessageSchemaVersion: 1,
		},
		Before: before,
		After:  after,
	}
}
