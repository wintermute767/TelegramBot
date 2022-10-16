package weather

import (
	"TelegramBot/src/telegramBot/processingMessage"
)

// Интерфейс для различных сайтов погода
type SiteWeather interface {
	GetWeatherAnswer(typeIncomingMassage string) *processingMessage.Answer
	GetForecast() *processingMessage.Answer
}
