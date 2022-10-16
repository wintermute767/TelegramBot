package databaseFnc

import (
	"TelegramBot/src/additionalFunc/databaseFnc/dao"
	"TelegramBot/src/telegramBot/processingMessage"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type DatabaseInterf interface {
	WriteCityToDB(database *dao.CityToDatabase) error
	GetCityFromDB(userID *int64) (string, error)
	WritLogToDB(*DataToDBLogs) error
}
type DataToDBLogs struct {
	Time, TypeMsg, UserId, IncomingMsg, OutgoingMsg string
}

func NewDataToDBLogs(update *tgbotapi.Update, answers ...*processingMessage.Answer) *DataToDBLogs {
	var time, typeMsg, userId, incomingMsg, outgoing string
	if update.Message != nil {
		time = strconv.Itoa(update.Message.Date)
		userId = strconv.FormatInt(update.Message.Chat.ID, 10)
		if update.Message.IsCommand() {
			typeMsg = "command"
			incomingMsg = "com:" + update.Message.Command() + "; arg:" + update.Message.CommandArguments()
		} else {
			typeMsg = "text"
			incomingMsg = update.Message.Text
		}
	}
	if update.CallbackQuery != nil {
		time = strconv.Itoa(update.CallbackQuery.Message.Date)
		userId = strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)
		typeMsg = "button"
		incomingMsg = "btn:" + update.CallbackQuery.Message.Text + "; data:" + update.CallbackQuery.Data
	}
	if answers != nil {
		outgoing = "foundSity: " + strings.Replace(answers[0].Query, "'", "", -1) + "; "
		for _, msg := range answers[0].AnswerMsg {
			if msg.ReplyMarkup != nil {
				outgoing += "buttons:" + strings.Replace(fmt.Sprint(msg.ReplyMarkup), "'", "", -1)
			}

		}
	}
	return &DataToDBLogs{Time: time, TypeMsg: typeMsg, UserId: userId, IncomingMsg: incomingMsg, OutgoingMsg: outgoing}
}
