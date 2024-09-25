package main

import (
	"database/sql"
	"rideshare-go/config"
	"rideshare-go/db"

	"github.com/sirupsen/logrus"
)

var Connection *sql.DB

func main() {
	log := logrus.New()

	config, err := config.LoadConfig("./config.json")
	if err != nil {
		log.WithError(err).Fatal("can't load config")
	}

	db, err := db.InitMySQL(log, config.MySQL)
	if err != nil {
		log.WithError(err).Fatal("can't initialize MySQL")
	}
	Connection = db
}
