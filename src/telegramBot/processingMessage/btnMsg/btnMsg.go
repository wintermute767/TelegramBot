package btnMsg

import (
	"TelegramBot/src/additionalFunc/weather/siteAccuweatherCom"
	"TelegramBot/src/telegramBot/processingMessage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Если нажмут кнопку

type BtnType struct {
	titleButton *string
	dataButton  *string
}

func NewBtnType(title *string, data *string) *BtnType {
	return &BtnType{titleButton: title, dataButton: data}
}

func (b BtnType) GetAnswerToMessage() *processingMessage.Answer {
	var answer processingMessage.Answer
	switch *b.titleButton {
	case "You need choose the exact name:":
		return siteAccuweatherCom.AccuWeather(*b.dataButton).GetWeatherAnswer("btn")
	case "Need a forecast?":
		return siteAccuweatherCom.AccuWeather(*b.dataButton).GetForecast()
	}
	return &answer

}
func (b BtnType) WriteInBot(Answer *processingMessage.Answer, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	//выводим на печать в боте
	//выводим на печать  название запрашиваемого города если есть
	if (*Answer).Query != "" {
		bot.Send(tgbotapi.NewMessage((*update).CallbackQuery.Message.Chat.ID, (*Answer).Query))
	}
	//выводим на печать в боте
	for _, msg := range (*Answer).AnswerMsg {
		msg.ChatID = (*update).CallbackQuery.Message.Chat.ID
		bot.Send(msg)
	}

}
