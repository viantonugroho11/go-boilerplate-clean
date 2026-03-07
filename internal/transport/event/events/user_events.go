package events

// UserCreatedEvent payload untuk event user created (Kafka JSON).
type UserCreatedEvent struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
