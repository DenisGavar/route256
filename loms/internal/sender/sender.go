package sender

import (
	"time"

	"github.com/Shopify/sarama"
)

type OrderMessage struct {
	OutboxKey int64
	Key       int64
	Message   string
	Timestamp time.Time
	Topic     string
}

type orderSender struct {
	producer sarama.SyncProducer
}

type Handler func(id string)

func NewOrderSender(producer sarama.SyncProducer, topic string) *orderSender {
	s := &orderSender{
		producer: producer,
	}

	return s
}
