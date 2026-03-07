package kafka

import (
	"context"
	"log"

	"go-boilerplate-clean/internal/config"
	kafkainfra "go-boilerplate-clean/internal/infrastructure/broker/kafka"
	"go-boilerplate-clean/internal/transport/event/events"
	userUsecase "go-boilerplate-clean/internal/usecase/users"

	"github.com/viantonugroho11/go-lib/kafka"
)

// UserCreatedHandler memproses UserCreatedEvent dan memanggil usecase.
type UserCreatedHandler struct {
	userService userUsecase.UserService
}

func NewUserCreatedHandler(userService userUsecase.UserService) *UserCreatedHandler {
	return &UserCreatedHandler{userService: userService}
}

func (h *UserCreatedHandler) Name() string { return "user-created" }

func (h *UserCreatedHandler) Handle(ctx context.Context, evt events.UserCreatedEvent, _ ...kafka.Header) kafka.Progress {
	if evt.ID == "" {
		return kafka.Progress{Status: kafka.ProgressDrop, Result: "id empty"}
	}
	user, err := h.userService.GetByID(ctx, evt.ID)
	if err != nil {
		log.Printf("user_consumer: GetByID %s: %v", evt.ID, err)
		p := kafka.Progress{Status: kafka.ProgressError, Result: err.Error()}
		p.SetError(err)
		return p
	}
	log.Printf("user_consumer: processed user id=%s name=%s", user.ID, user.Name)
	return kafka.Progress{Status: kafka.ProgressSuccess, Result: "ok"}
}

// RunUser menjalankan consumer user-created.
func RunUser(ctx context.Context, cfg *config.Configuration, userService userUsecase.UserService) (kafka.Consumer, error) {
	return kafkainfra.RunWithConfig[events.UserCreatedEvent](ctx, cfg, cfg.Kafka.GroupID, cfg.Kafka.Topic, NewUserCreatedHandler(userService))
}
