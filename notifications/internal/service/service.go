package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"route256/libs/logger"
	"route256/notifications/internal/kafka"
	"route256/notifications/internal/pkg/notifications"
	"route256/notifications/internal/telegram"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type Repository interface {
	GetHistory(ctx context.Context, req *notifications.GetHistoryRequest) (*notifications.GetHistoryResponse, error)
	SaveNotification(ctx context.Context, orderId, userId int64, status string) error
}

type Cache interface {
	Add(key string, value any)
	Get(key string) any
}

type Service struct {
	notifications.UnimplementedNotificationsServer
	TransactionManager
	Repository
	Cache
	Brokers []string
	Bot     *telegram.Bot
}

func NewNotificationService(provider TransactionManager, repo Repository, cache Cache, brokers []string, bot *telegram.Bot) *Service {
	return &Service{
		TransactionManager: provider,
		Repository:         repo,
		Cache:              cache,
		Brokers:            brokers,
		Bot:                bot,
	}
}

/*
Позволяет получить историю уведомлений пользователя за указанный промежуток времени.
Пример:

	{
	    "user_id": "1",
	    "start_time": "2023-07-15 06:54:11+03",
	    "end_time": "2023-07-15 07:23:43+03"
	}

start_time и end_time строго 22 знака и именно в таком формате
+03 - это время относительно UTC, в бд время хранится в timestamptz
*/
func (s *Service) GetHistory(ctx context.Context, req *notifications.GetHistoryRequest) (*notifications.GetHistoryResponse, error) {
	start := time.Now()
	logger.Infof("%+v\n", req)
	err := req.ValidateAll()
	if err != nil {
		logger.Infoln("ERROR", err)
		return nil, err
	}
	key := fmt.Sprintf("%d--%s--%s", req.UserId, req.StartTime, req.EndTime)
	data := s.Cache.Get(key)
	result, ok := data.(*notifications.GetHistoryResponse)
	if ok {
		logger.Infoln(time.Since(start))
		return result, nil
	}
	result, err = s.Repository.GetHistory(ctx, req)
	logger.Infoln(time.Since(start))
	if err != nil {
		return nil, err
	}
	s.Cache.Add(key, result)
	return result, nil
}

func (s *Service) ConsumerGroupStatuses(brokers []string, bot *telegram.Bot, repo Repository) {
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

	consumer := kafka.NewConsumerGroup(bot, repo)
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
