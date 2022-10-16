package urlWeatherFunc

import "regexp"

// Функция для уменьшения длины Url в кнопке из-за ограничений
func CutUrlBtnChoose(str string) string {
	reg := regexp.MustCompile(`=.*&`)
	str = reg.FindString(str)
	return str
}
func CompleteUrlBtnChoose(str string) string {
	url := "https://www.accuweather.com/web-api/three-day-redirect?key" + str + "target="
	return url
}

// Функция для уменьшения длины Url
func CutUrlBtnForecast(str string) string {
	reg := regexp.MustCompile(`\/..\/..\/(.*)\/daily-weather-forecast.*`)
	strGroup := reg.FindSubmatch([]byte(str))
	return string(strGroup[1])
}

// Получение полной версии Url
func CompleteUrlBtnForecast(str string) string {
	reg := regexp.MustCompile(`\w+/(.*)`)
	strGroup := reg.FindSubmatch([]byte(str))
	urlForecast := "/en/ru/" + str + "/daily-weather-forecast/" + string(strGroup[1])
	return urlForecast
}
