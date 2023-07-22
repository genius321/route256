package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"route256/libs/logger"
	"route256/notifications/internal/telegram"
	"strconv"

	"github.com/Shopify/sarama"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type statusMessage struct {
	OrderId    int64
	UserId     int64
	StatusName string
}

type Repository interface {
	SaveNotification(ctx context.Context, orderId, userId int64, status string) error
}

type ConsumerGroup struct {
	Repository
	bot   *telegram.Bot
	ready chan bool
}

func NewConsumerGroup(bot *telegram.Bot, repo Repository) ConsumerGroup {
	return ConsumerGroup{
		Repository: repo,
		bot:        bot,
		ready:      make(chan bool),
	}
}

func (consumer *ConsumerGroup) Ready() <-chan bool {
	return consumer.ready
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (consumer *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)

	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (consumer *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор пока сессия не завершилась
func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			sm := statusMessage{}
			err := json.Unmarshal(message.Value, &sm)
			if err != nil {
				logger.Error("Consumer group error", err)
			}

			response := fmt.Sprintf("orderId:%d userId: %d status:%s\n", sm.OrderId, sm.UserId, sm.StatusName)
			chatID, err := strconv.ParseInt(os.Getenv("TG_CHAT"), 10, 64)
			ctx := context.Background()
			if err != nil {
				logger.Errorf(ctx, "ParseInt", "retrieve chatID: %w", err)
			}
			msg := tgbotapi.NewMessage(chatID, response)
			consumer.bot.SendMessage(&msg)

			err = consumer.SaveNotification(ctx, sm.OrderId, sm.UserId, sm.StatusName)
			if err != nil {
				logger.Errorf(ctx, "SaveNotification", "SaveNotification: %w", err)
			}
			logger.Infof("Message claimed: value = %v, timestamp = %v, topic = %s",
				sm,
				message.Timestamp,
				message.Topic,
			)

			// коммит сообщения "руками"
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
