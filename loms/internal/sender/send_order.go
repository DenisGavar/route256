package sender

import (
	"context"
	"fmt"
	"route256/libs/logger"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
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

	logger.Debug("sending order", zap.Int64("key", orderMessage.Key), zap.Int32("partition", partition), zap.Int64("offset", offset))

	return nil
}
