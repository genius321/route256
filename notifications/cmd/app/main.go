package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"

	"route256/notifications/internal/service"
	"route256/notifications/internal/telegram"
)

func main() {
	// создаёт бота
	bot, err := tgbotapi.NewBotAPI("6376011717:AAFrInwAk9DH2NjwbRXuBCslp9D6QVHsPew")
	if err != nil {
		log.Panic(err)
	}

	telegramBot := telegram.NewBot(bot)

	// список брокеров
	var brokers = []string{
		"kafka1:29091",
		"kafka2:29092",
		"kafka3:29093",
	}

	// создаёт сервис нотификаций
	s := service.NewService(brokers, telegramBot)
	s.ConsumerGroupStatuses(brokers, telegramBot)
}
