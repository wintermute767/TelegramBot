package additionalWetherFunc

import (
	"TelegramBot/src/additionalFunc/domSite"
	"TelegramBot/src/additionalFunc/weather/siteAccuweatherCom/urlWeatherFunc"
	"TelegramBot/src/telegramBot/processingMessage"
	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/url"
	"regexp"
	"strings"
)

func GetUrl(typeIncomingMassage string, searchSity string) string {
	var urlLink string
	//в зависимости от типа входящего сообщения будем по разному запрашивать погоду
	switch typeIncomingMassage {
	case "text":
		re := regexp.MustCompile("[A-Za-z]")
		if re.Match([]byte(strings.ToLower(searchSity))) {
			return "https://www.accuweather.com/en/search-locations?query=" + url.QueryEscape(searchSity)
		} else {
			return "https://www.accuweather.com/ru/search-locations?query=" + url.QueryEscape(searchSity)
		}
	case "btn":
		return "https://www.accuweather.com/web-api/three-day-redirect?key" + searchSity + "target="
	}
	return urlLink
}

func GetTemp(dom *goquery.Document) (string, bool) {
	var ok bool = false
	var query string
	dom.Find("div.cur-con-weather-card__body>div>div.forecast-container>div.temp-container>div.temp").Each(func(i int, selection *goquery.Selection) {
		query = "Temperature in the selected city: " + selection.Text()
		ok = true
	})
	return query, ok
}

func CheckingResult(dom *goquery.Document) (string, bool) {
	var ok bool = false
	var query string
	dom.Find("div.no-results-text").Each(func(i int, selection *goquery.Selection) {
		query = selection.Text()
		ok = true
	})
	return query, ok
}

func CompilationOfList(searchCity string, dom *goquery.Document) (processingMessage.Answer, bool) {
	var ok bool = false
	listRezult := make(map[string]string)
	var urlExactMatch string
	dom.Find("div.locations-list>a").Each(func(i int, selection *goquery.Selection) {
		link, _ := selection.Attr("href")
		//удалить повторяющиеся строки в кнопках
		//при повторах выбираем последнее
		reg := regexp.MustCompile("\\s*\\(.*\\)")
		res := reg.ReplaceAllString(strings.TrimSpace(selection.Text()), "")
		listRezult[res] = link
		ok = true
		//проверяем есть ли полное совпадение
		re := regexp.MustCompile("^" + strings.ToLower(searchCity) + ",\\s" + strings.ToLower(searchCity))
		if re.Match([]byte(strings.ToLower(res))) {
			urlExactMatch = string(link)
		}
	})
	var answer processingMessage.Answer
	//Выбираем если есть точное совпадение то еще раз проводим поиск температуры, если нет то выводим список кнопок
	if urlExactMatch != "" {
		urlLink := "https://www.accuweather.com" + urlExactMatch
		dom = domSite.GetDomSite(urlLink)
		temp, _ := GetTemp(dom)
		// Добавим температуру, сообщение с деталями, кнопку с сылкой на прогоз погоды
		answer.AnswerMsg = append(answer.AnswerMsg, tgbotapi.NewMessage(0, temp), tgbotapi.NewMessage(0, GetDetails(dom)), GetBtnForecast(dom))
		//Добавим название города для дальнейшей записи в БД в структуру ответа
		dom.Find("div.header-inner>a.header-city-link>h1.header-loc").Each(func(i int, selection *goquery.Selection) {
			answer.Query = selection.Text()
		})
		return answer, ok
	} else {
		answer.AnswerMsg = append(answer.AnswerMsg, createBtnChoose(listRezult))
	}
	return answer, ok
}

func createBtnChoose(listRezult map[string]string) tgbotapi.MessageConfig {
	btnList := tgbotapi.NewMessage(0, "")
	btnList.Text = "You need choose the exact name:"
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for key, value := range listRezult {
		var row []tgbotapi.InlineKeyboardButton
		//необходимо уменьшить длину ссылки value из-за ограничекния в 64Bt поэтому  urlToBtn.CutUrlBtnChoose
		btn := tgbotapi.NewInlineKeyboardButtonData(key, urlWeatherFunc.CutUrlBtnChoose(value))
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	btnList.ReplyMarkup = keyboard
	return btnList
}

// Получить сообщение с деталями погоды
func GetDetails(dom *goquery.Document) string {
	//Найдем ссылку на подробности
	var urlDetails string
	dom.Find("a[data-gaid=hourly]").Each(func(i int, selection *goquery.Selection) {
		urlDetails, _ = selection.Attr("href")
	})
	//ищем на новой странице детали
	var textDetail []string
	dom = domSite.GetDomSite("https://www.accuweather.com" + urlDetails)
	dom.Find("div#hourlyCard0>div.accordion-item-content>div.hourly-card-nfl-content>div.hourly-content-container>div.panel.left>p>span").Each(func(i int, selection *goquery.Selection) {
		textDetail = append(textDetail, selection.Text())
	})
	if len(textDetail) != 0 {
		details := "Wind " + textDetail[0] + "\nWind Gusts: " + textDetail[1] + "\nHumidity:  " + textDetail[2] + "\nIndoor Humidity:  " + textDetail[3] + "\nAir Quality:  " + textDetail[4]
		return details
	}
	return "There are no details"
}

func GetBtnForecast(dom *goquery.Document) tgbotapi.MessageConfig {
	//получаем ссылку на новую страницу
	var urlForecast string
	dom.Find("a[data-gaid=daily]").Each(func(i int, selection *goquery.Selection) {
		urlForecast, _ = selection.Attr("href")
	})
	btnForecast := tgbotapi.NewMessage(0, "Need a forecast?")
	keyboardDetail := tgbotapi.InlineKeyboardMarkup{}
	//необходимо уменьшить длину ссылки  из-за ограничекния в 64Bt поэтому  urlToBtn.CutUrlBtnForecast
	rowDetail := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Forecast", urlWeatherFunc.CutUrlBtnForecast(urlForecast))}
	keyboardDetail.InlineKeyboard = append(keyboardDetail.InlineKeyboard, rowDetail)
	btnForecast.ReplyMarkup = keyboardDetail
	return btnForecast
}
