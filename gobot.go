package main

import (
	"TelegramBot/src/additionalFunc/databaseFnc"
	"TelegramBot/src/additionalFunc/databaseFnc/dao"
	"TelegramBot/src/additionalFunc/databaseFnc/postgreSQL"
	"TelegramBot/src/telegramBot"
	"TelegramBot/src/telegramBot/processingMessage"
	"TelegramBot/src/telegramBot/processingMessage/btnMsg"
	"TelegramBot/src/telegramBot/processingMessage/comMsg"
	"TelegramBot/src/telegramBot/processingMessage/txtMsg"
	"fmt"
	"os"
)

func main() {
	token := os.Getenv("TOKEN")
	bot, updates := telegramBot.StartTBot(&token)

	var databaseInter databaseFnc.DatabaseInterf = postgreSQL.NewPostgreSQL()

	for update := range updates {
		fmt.Println(bot, update)
		var incomingMassage processingMessage.NewMassage

		//Вынуждены применять усоривие так как не все виды сообщений имеют собственный тип
		//задаем тип каждому виду сообщений
		if update.Message != nil {
			if update.Message.IsCommand() {
				com := update.Message.Command()
				arg := update.Message.CommandArguments()
				//если команда w то аргументом будет значение из базы данных
				if com == "w" {
					//для команды /w - ищем аргументы в БД
					argFromDB, err := databaseInter.GetCityFromDB(&update.Message.Chat.ID)
					if err != nil {
						panic(err)
					} else {
						arg = argFromDB
					}
				}
				incomingMassage = comMsg.NewCommandType(&com, &arg)
			} else {
				incomingMassage = txtMsg.NewTextType(&update.Message.Text)
			}
		}
		if update.CallbackQuery != nil {
			incomingMassage = btnMsg.NewBtnType(&update.CallbackQuery.Message.Text, &update.CallbackQuery.Data)
		}

		//Запишим входящее значение в лог
		go databaseInter.WritLogToDB(databaseFnc.NewDataToDBLogs(&update))
		//создаем канал для результатов ответа
		answerChan := make(chan *processingMessage.Answer)
		go incomingMassage.GetAnswerToMessage(answerChan)
		//ждем когда получим ответ из канала
		answers := <-answerChan
		go incomingMassage.WriteInBot(answers, bot, &update)
		//Запишем в базу данных город который успешно нашли в этот раз для последующих поисков по команде
		go databaseInter.WriteCityToDB(dao.NewCityToDatabase(&update, &(answers).Query))
		//Запишим ответ в лог
		go databaseInter.WritLogToDB(databaseFnc.NewDataToDBLogs(&update, answers))

	}

}
