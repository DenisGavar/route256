package sender

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

func (s *orderSender) SendOrder(ctx context.Context, orderMessage *OrderMessage) error {
	msg := &sarama.ProducerMessage{
		Topic:     orderMessage.Topic,
		Partition: -1,
		Value:     sarama.StringEncoder(orderMessage.Message),
		Key:       sarama.StringEncoder(fmt.Sprint(orderMessage.Key)),
		Timestamp: orderMessage.Timestamp,
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("key: %d, partition: %d, offset: %d", orderMessage.Key, partition, offset)
	return nil
}
