package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"order-service/internal/event"
)

type OrderProducer struct {
	writer *kafka.Writer
}

func NewOrderProducer(brokers []string, topic string) *OrderProducer {
	return &OrderProducer{writer: &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}}
}

func (p *OrderProducer) PublishOrderCreated(ctx context.Context, evt event.OrderCreated) error {
	evtJson, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(evt.UserID),
		Value: evtJson,
	})
}
