package main

import (
	"context"
	"os"
	"os/signal"
	"route256/libs/kafka"
	"route256/libs/logger"
	configServices "route256/notifications/internal/config"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

func main() {
	logger.Init()

	keepRunning := true
	logger.Info("Starting a new Sarama consumer")

	err := configServices.Init()
	if err != nil {
		logger.Fatal("config init", zap.Error(err))
	}

	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	switch configServices.ConfigData.Services.Kafka.BalanceStrategy {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:
		logger.Fatal("Unrecognized consumer group partition assignor", zap.String("Balance strategy", configServices.ConfigData.Services.Kafka.BalanceStrategy))
	}

	consumer := kafka.NewConsumerGroup()

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(configServices.ConfigData.Services.Kafka.Brokers, configServices.ConfigData.Services.Kafka.GroupName, config)
	if err != nil {
		logger.Fatal("Error creating consumer group client", zap.Error(err))
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{configServices.ConfigData.Services.Kafka.TopicForOrders}, &consumer); err != nil {
				logger.Fatal("Error from consumer", zap.Error(err))
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.Ready()
	logger.Info("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.Signal(0xa)) // syscall.SIGUSR1

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			logger.Info("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			logger.Info("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}

	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		logger.Fatal("Error closing client", zap.Error(err))
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		logger.Info("Resuming consumption")
	} else {
		client.PauseAll()
		logger.Info("Pausing consumption")
	}

	*isPaused = !*isPaused
}
