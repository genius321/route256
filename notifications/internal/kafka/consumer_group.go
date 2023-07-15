package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/libs/logger"
	"route256/notifications/internal/telegram"

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
			msg := tgbotapi.NewMessage(457312730, response)
			consumer.bot.SendMessage(&msg)

			err = consumer.SaveNotification(context.Background(), sm.OrderId, sm.UserId, sm.StatusName)
			if err != nil {
				logger.Error("Consumer group SaveNotification error", err)
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
