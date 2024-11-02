package main

import (
	"main/api"
	"main/auth"
	"main/core"
	"main/db"
	"net/http"
)

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

	r := api.CreateRouter(db, config.Server.AuthSecret)

	googleAuthModule := auth.NewGoogleAuthModule(config.GoogleAuth, db, config.Server.AuthSecret)
	googleAuthModule.ApplyRoutes(r)

	port := config.Server.Port
	log.WithField("port", port).Info("starting server")

	log.Fatal(http.ListenAndServe(":"+port, r))
}
