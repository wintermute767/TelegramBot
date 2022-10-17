package processingMessage

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Answer struct {
	Query     string
	AnswerMsg []tgbotapi.MessageConfig
}

type NewMassage interface {
	GetAnswerToMessage(chan *Answer)
	WriteInBot(*Answer, *tgbotapi.BotAPI, *tgbotapi.Update)
}
