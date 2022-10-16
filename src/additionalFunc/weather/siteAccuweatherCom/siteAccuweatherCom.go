package siteAccuweatherCom

import (
	"TelegramBot/src/additionalFunc/domSite"
	"TelegramBot/src/additionalFunc/weather/siteAccuweatherCom/additionalWetherFunc"
	"TelegramBot/src/additionalFunc/weather/siteAccuweatherCom/urlWeatherFunc"

	"TelegramBot/src/telegramBot/processingMessage"
	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type AccuWeather string

func (a AccuWeather) GetWeatherAnswer(typeIncomingMassage string) *processingMessage.Answer {
	var answer processingMessage.Answer

	urlLink := additionalWetherFunc.GetUrl(typeIncomingMassage, string(a))
	//Ищем температуру
	dom := domSite.GetDomSite(urlLink)
	temp, ok := additionalWetherFunc.GetTemp(dom)
	if ok {
		// Добавим температуру, сообщение с деталями, кнопку с сылкой на прогоз погоды
		answer.AnswerMsg = append(answer.AnswerMsg, tgbotapi.NewMessage(0, temp), tgbotapi.NewMessage(0, additionalWetherFunc.GetDetails(dom)), additionalWetherFunc.GetBtnForecast(dom))
		//Добавим название города для дальнейшей записи в БД в структуру ответа
		dom.Find("div.header-inner>a.header-city-link>h1.header-loc").Each(func(i int, selection *goquery.Selection) {
			answer.Query = selection.Text()
		})
		return &answer
	}
	//Если сразу не нашли, проверям сообщение с сайта о результате поиска
	result, ok := additionalWetherFunc.CheckingResult(dom)
	if ok {
		answer.AnswerMsg = append(answer.AnswerMsg, tgbotapi.NewMessage(0, result))
		return &answer
	}
	//Если резльтов поиска нет проверяем может есть список возможных результатов
	answerListBtn, ok := additionalWetherFunc.CompilationOfList(string(a), dom)
	if ok {
		return &answerListBtn
	}
	//Если совсем что то пошло не так и ничего не сработало

	answer.AnswerMsg = append(answer.AnswerMsg, tgbotapi.NewMessage(0, "Error! Something went wrong("))
	return &answer
}

func (a AccuWeather) GetForecast() *processingMessage.Answer {
	var answer processingMessage.Answer
	//Переходим на новую страницу
	var day, date, max, low, forecast string
	dom := domSite.GetDomSite("https://www.accuweather.com" + urlWeatherFunc.CompleteUrlBtnForecast(string(a)))
	//Добавим название города для дальнейшей записи в БД в структуру ответа
	dom.Find("div.header-inner>a.header-city-link>h1.header-loc").Each(func(i int, selection *goquery.Selection) {
		answer.Query = selection.Text()
	})

	dom.Find("div.page-content.content-module>div.daily-wrapper>a>div.info").Each(func(i int, selection *goquery.Selection) {

		selection.Find("h2>span.module-header.dow.date").Each(func(i int, selection *goquery.Selection) {
			day = strings.TrimSpace(selection.Text())
		})
		selection.Find("h2>span.module-header.sub.date").Each(func(i int, selection *goquery.Selection) {
			date = strings.TrimSpace(selection.Text())
		})
		selection.Find("div>span.high").Each(func(i int, selection *goquery.Selection) {
			max = strings.TrimSpace(selection.Text())
		})
		selection.Find("div>span.low").Each(func(i int, selection *goquery.Selection) {
			low = strings.TrimSpace(selection.Text())
		})
		forecast += day + " (date " + date + ") temperature (min\\max):    " + max + "\\" + low[1:] + "\n"
	})

	answer.AnswerMsg = append(answer.AnswerMsg, tgbotapi.NewMessage(0, forecast))
	return &answer
}
