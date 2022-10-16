package postgreSQL

import (
	"TelegramBot/src/additionalFunc/databaseFnc"
	"TelegramBot/src/additionalFunc/databaseFnc/dao"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strconv"
	"strings"
)

type postgreSQL struct {
	connectToDB *pgx.Conn
}

func stopPanic() {
	recover()
}

func (p *postgreSQL) WriteCityToDB(data *dao.CityToDatabase) error {
	var err error
	if data.City != "" {
		var insertStmt string = "INSERT INTO lastsearch (userid, city) VALUES ('" + strconv.FormatInt(data.UserID, 10) + "', '" + strings.Replace(data.City, "'", "", -1) + "') ON CONFLICT (userid) DO UPDATE SET city = EXCLUDED.city;"
		_, err = p.connectToDB.Exec(context.Background(), insertStmt)
	}
	if err != nil {
		panic(err)
	}
	return err
}

func (p *postgreSQL) GetCityFromDB(userID *int64) (string, error) {
	var city string
	defer stopPanic()
	err := p.connectToDB.QueryRow(context.Background(), "SELECT city  FROM lastsearch WHERE userid= '"+strconv.FormatInt(*userID, 10)+"'").Scan(&city)
	if err != nil {
		panic(err)
	}
	return city, err
}

func NewPostgreSQL() *postgreSQL {
	databaseUrl := "postgres://postgres:postgrespw@host.docker.internal:55000/telegram"
	connectToDB, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		panic(err)
	}
	return &postgreSQL{connectToDB: connectToDB}
}
func (p *postgreSQL) WritLogToDB(dataToDBLogs *databaseFnc.DataToDBLogs) error {
	var err error

	var insertStmt string = "INSERT INTO logs (time, type,userid, incomingmsg, outgoing ) " +
		"VALUES (" +
		"to_timestamp(" + dataToDBLogs.Time + ")," +
		"'" + dataToDBLogs.TypeMsg + "'," +
		"'" + dataToDBLogs.UserId + "'," +
		"'" + dataToDBLogs.IncomingMsg + "'," +
		"'" + dataToDBLogs.OutgoingMsg + "') ;"

	fmt.Println(insertStmt)
	_, err = p.connectToDB.Exec(context.Background(), insertStmt)

	if err != nil {
		panic(err)
	}

	return err
}
