package txtMsg

import (
	"TelegramBot/src/additionalFunc/weather/siteAccuweatherCom"
	"TelegramBot/src/telegramBot/processingMessage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TextType struct {
	search *string
}

func NewTextType(txt *string) *TextType {
	return &TextType{search: txt}
}

func (t *TextType) GetAnswerToMessage(answersChan chan *processingMessage.Answer) {
	answersChan <- siteAccuweatherCom.AccuWeather(*t.search).GetWeatherAnswer("text")
}

func (t *TextType) WriteInBot(Answer *processingMessage.Answer, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	//выводим на печать в боте
	//выводим на печать  название запрашиваемого города если есть
	if (*Answer).Query != "" {
		bot.Send(tgbotapi.NewMessage((*update).Message.Chat.ID, (*Answer).Query))
	}
	//выодим ответ
	for _, msg := range (*Answer).AnswerMsg {
		msg.ChatID = (*update).Message.Chat.ID
		bot.Send(msg)
	}
}
