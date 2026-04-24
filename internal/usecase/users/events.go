package users

import (
	"context"
	userEntity "go-boilerplate-clean/internal/entity/users"
)

// UserEventPublisher interface untuk publish event user (mis. ke Kafka).
// Implementasi bisa menggunakan go-lib/kafka Producer.
type UserEventPublisher interface {
	PublishUser(ctx context.Context, user userEntity.User) error
}
