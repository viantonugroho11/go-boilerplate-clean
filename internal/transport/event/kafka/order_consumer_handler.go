package kafka

import (
	"context"
	"log"

	"go-boilerplate-clean/internal/config"
	kafkainfra "go-boilerplate-clean/internal/infrastructure/broker/kafka"
	"go-boilerplate-clean/internal/transport/event/events"

	"github.com/viantonugroho11/go-lib/kafka"
)

// OrderCreatedHandler memproses OrderCreatedEvent (bisa inject order usecase nanti).
type OrderCreatedHandler struct{}

func NewOrderCreatedHandler() *OrderCreatedHandler { return &OrderCreatedHandler{} }

func (h *OrderCreatedHandler) Name() string { return "order-created" }

func (h *OrderCreatedHandler) Handle(ctx context.Context, evt events.OrderCreatedEvent, _ ...kafka.Header) kafka.Progress {
	if evt.ID == "" {
		return kafka.Progress{Status: kafka.ProgressDrop, Result: "order id empty"}
	}
	log.Printf("order_consumer: order id=%s user_id=%s amount=%.2f", evt.ID, evt.UserID, evt.Amount)
	return kafka.Progress{Status: kafka.ProgressSuccess, Result: "ok"}
}

// RunOrder menjalankan consumer order-created.
func RunOrder(ctx context.Context, cfg *config.Configuration) (kafka.Consumer, error) {
	groupID, topic := cfg.Kafka.GroupIDOrders, cfg.Kafka.TopicOrders
	if groupID == "" {
		groupID = "order-consumer-group"
	}
	if topic == "" {
		topic = "orders"
	}
	return kafkainfra.RunWithConfig[events.OrderCreatedEvent](ctx, cfg, groupID, topic, NewOrderCreatedHandler())
}
