package status

import (
	"encoding/json"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/kafka"
	orderModels "route256/loms/internal/models/order"

	"github.com/Shopify/sarama"
)

type statusMessage struct {
	OrderId    int64
	UserId     int64
	StatusName string
}

type KafkaSender struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaSender(producer *kafka.Producer, topic string) *KafkaSender {
	return &KafkaSender{
		producer,
		topic,
	}
}

func (s *KafkaSender) SendMessage(orderId orderModels.OrderId, userId orderModels.User, stataus orderModels.Status) error {
	message := statusMessage{OrderId: int64(orderId), UserId: int64(userId), StatusName: string(stataus)}
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		logger.Infoln("Send message marshal error", err)
		return err
	}

	partition, offset, err := s.producer.SendSyncMessage(kafkaMsg)

	if err != nil {
		logger.Infoln("Send message connector error", err)
		return err
	}

	logger.Infoln("Partition: ", partition, " Offset: ", offset, " AnswerID:", message.OrderId)
	return nil
}

func (s *KafkaSender) buildMessage(message statusMessage) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)

	if err != nil {
		logger.Infoln("Send message marshal error", err)
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(message.OrderId)),
	}, nil
}
