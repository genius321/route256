package main

import (
	"flag"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"route256/libs/logger"
	"route256/notifications/internal/service"
	"route256/notifications/internal/telegram"
)

var (
	environment = flag.String("environment", "DEVELOPMENT", "environment: [DEVELOPMENT, PRODUCTION]")
)

func main() {
	flag.Parse()
	// Init logger for environment
	logger.SetLoggerByEnvironment(*environment)

	// создаёт бота
	bot, err := tgbotapi.NewBotAPI("6376011717:AAFrInwAk9DH2NjwbRXuBCslp9D6QVHsPew")
	if err != nil {
		logger.Panic(err)
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
