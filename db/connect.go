package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

var Conn *sqlx.DB

func Init() {
	var err error
	var host = os.Getenv("POSTGRES_HOST")
	var username = os.Getenv("POSTGRES_USER")
	var password = os.Getenv("POSTGRES_PASSWORD")
	var dbName = os.Getenv("POSTGRES_DB")
	if password != "" {
		username += ":" + password
	}
	var connection = "postgres://" + username + "@" + host + "/" + dbName + "?sslmode=disable"

	Conn, err = sqlx.Connect("postgres", connection)
	if err != nil {
		logrus.Fatalln(err)
	}
}
