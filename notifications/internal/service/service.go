package service

import (
	"context"
	"os"
	"os/signal"
	"route256/libs/logger"
	"route256/notifications/internal/kafka"
	"route256/notifications/internal/telegram"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

type Service struct {
	Brokers []string
	Bot     *telegram.Bot
}

func NewService(brokers []string, bot *telegram.Bot) *Service {
	return &Service{Brokers: brokers, Bot: bot}
}

func (s *Service) ConsumerGroupStatuses(brokers []string, bot *telegram.Bot) {
	keepRunning := true
	logger.Infoln("Starting a new Sarama consumer")

	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Используется, если ваш offset "уехал" далеко и нужно пропустить невалидные сдвиги
	config.Consumer.Group.ResetInvalidOffsets = true

	// Сердцебиение консьюмера
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	// Таймаут сессии
	config.Consumer.Group.Session.Timeout = 60 * time.Second

	// Таймаут ребалансировки
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second

	const BalanceStrategy = "roundrobin"
	switch BalanceStrategy {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:
		logger.Panicf("Unrecognized consumer group partition assignor: %s", BalanceStrategy)
	}

	consumer := kafka.NewConsumerGroup(bot)
	group := "notifications"

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		logger.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{"statuses"}, &consumer); err != nil {
				logger.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.Ready() // Await till the consumer has been set up
	logger.Infoln("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			logger.Infoln("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			logger.Infoln("terminating: via signal")
			keepRunning = false
		}
	}

	cancel()
	wg.Wait()

	if err = client.Close(); err != nil {
		logger.Panicf("Error closing client: %v", err)
	}
}
