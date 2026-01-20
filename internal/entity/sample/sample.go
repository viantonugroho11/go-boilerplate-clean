package sample

import "time"

type Sample struct {
	ID        string     `json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

const (
	SampleStatusOpen   = "open"
	SampleStatusOnHold = "on_hold"
	SampleStatusClosed = "closed"
)
