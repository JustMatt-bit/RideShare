package main

import (
	"database/sql"
	"net/http"
)

type any interface{}

var Connection *sql.DB

func main() {
	log := setupLogging()

	config, err := loadConfig("./config.json")
	if err != nil {
		log.WithError(err).Fatal("can't load config")
	}

	db, err := initMySQL(log, config.MySQL)
	if err != nil {
		log.WithError(err).Fatal("can't initialize MySQL")
	}
	defer db.Close()
	Connection = db

	r := createRouter()

	port := config.Server.Port
	log.WithField("port", port).Info("starting server")

	log.Fatal(http.ListenAndServe(":"+port, r))
}
