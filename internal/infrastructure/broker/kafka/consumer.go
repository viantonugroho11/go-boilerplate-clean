package kafka

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type MessageHandler func(ctx context.Context, msg *sarama.ConsumerMessage) error

type Consumer struct {
	group   sarama.ConsumerGroup
	topic   string
	handler MessageHandler

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewConsumer(brokers []string, groupID string, topic string, handler MessageHandler) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Version = sarama.V2_8_0_0

	group, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		group:   group,
		topic:   topic,
		handler: handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if err := c.group.Consume(ctx, []string{c.topic}, &cgHandler{handler: c.handler}); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Printf("kafka consume error: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
}

func (c *Consumer) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	return c.group.Close()
}

type cgHandler struct {
	handler MessageHandler
}

func (h *cgHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *cgHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *cgHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		_ = h.handler(sess.Context(), msg)
		sess.MarkMessage(msg, "")
	}
	return nil
}
