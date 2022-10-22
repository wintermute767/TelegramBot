package postgreSQL

import (
	"TelegramBot/src/additionalFunc/databaseFnc"
	"TelegramBot/src/additionalFunc/databaseFnc/dao"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strconv"
	"strings"
)

type postgreSQL struct {
	connectToDB *pgxpool.Pool
}

func NewPostgreSQL() *postgreSQL {
	databaseUrl := os.Getenv("POSTGRES_CONNECT")
	connectToDB, err := pgxpool.New(context.Background(), databaseUrl)

	if err != nil {
		panic(err)
	}
	return &postgreSQL{connectToDB: connectToDB}
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

	err := p.connectToDB.QueryRow(context.Background(), "SELECT city  FROM lastsearch WHERE userid= '"+strconv.FormatInt(*userID, 10)+"'").Scan(&city)
	if err != nil {
		city = ""
	}

	return city, err
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

	_, err = p.connectToDB.Exec(context.Background(), insertStmt)
	if err != nil {
		panic(err)
	}
	return err
}
