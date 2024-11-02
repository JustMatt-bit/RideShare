package main

import (
	"database/sql"
	"main/api"
	"main/auth"
	"main/core"
	"main/db"
	"net/http"

	"golang.org/x/oauth2"
)

type any interface{}

var Connection *sql.DB
var GoogleAuth oauth2.Config

func main() {
	log := core.SetupLogging()

	config, err := core.LoadConfig("./config.json")
	if err != nil {
		log.WithError(err).Fatal("can't load config")
	}

	db, err := db.InitMySQL(log, config.MySQL)
	if err != nil {
		log.WithError(err).Fatal("can't initialize MySQL")
	}
	defer db.Close()
	Connection = db

	r := api.CreateRouter(db)

	googleAuthModule := auth.NewGoogleAuthModule(config.GoogleAuth, db)
	googleAuthModule.ApplyRoutes(r)

	port := config.Server.Port
	log.WithField("port", port).Info("starting server")

	log.Fatal(http.ListenAndServe(":"+port, r))
}
