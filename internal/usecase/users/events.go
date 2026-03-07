package users

import "context"

// UserEventPublisher interface untuk publish event user (mis. ke Kafka).
// Implementasi bisa menggunakan go-lib/kafka Producer.
type UserEventPublisher interface {
	PublishUserCreated(ctx context.Context, id, name, email string) error
}
