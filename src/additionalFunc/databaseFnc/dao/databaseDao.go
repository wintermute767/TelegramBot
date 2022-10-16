package dao

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CityToDatabase struct {
	UserID int64
	City   string
}

func NewCityToDatabase(update *tgbotapi.Update, city *string) *CityToDatabase {
	var ID int64
	if update.Message != nil {
		ID = update.Message.Chat.ID
	} else {
		ID = update.CallbackQuery.Message.Chat.ID
	}
	return &CityToDatabase{UserID: ID, City: *city}
}
