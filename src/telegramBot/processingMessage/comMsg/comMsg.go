package comMsg

import (
	"TelegramBot/src/additionalFunc/weather/siteAccuweatherCom"
	"TelegramBot/src/telegramBot/processingMessage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Если пришлют команду

type CommandType struct {
	commandName *string
	commandArg  *string
}

func NewCommandType(command *string, arg *string) *CommandType {
	return &CommandType{commandName: command, commandArg: arg}
}

func (c CommandType) GetAnswerToMessage() *processingMessage.Answer {
	var answer processingMessage.Answer

	switch *c.commandName {
	case "w":
		return siteAccuweatherCom.AccuWeather(*c.commandArg).GetWeatherAnswer("text")
	case "h":
		answer.AnswerMsg = append(answer.AnswerMsg, tgbotapi.NewMessage(0, "Enter the command /w to get the result of the old success request"))
	}

	return &answer

}
func (c CommandType) WriteInBot(Answer *processingMessage.Answer, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	//выводим на печать в боте
	//выводим на печать  название запрашиваемого города если есть
	if (*Answer).Query != "" {
		bot.Send(tgbotapi.NewMessage((*update).Message.Chat.ID, (*Answer).Query))
	}
	//выводим на печать ответ
	for _, msg := range (*Answer).AnswerMsg {
		msg.ChatID = (*update).Message.Chat.ID
		bot.Send(msg)
	}
}
