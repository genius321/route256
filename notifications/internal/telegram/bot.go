package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{bot: bot}
}

func (b *Bot) SendMessage(msg *tgbotapi.MessageConfig) {
	b.bot.Send(msg)
}
