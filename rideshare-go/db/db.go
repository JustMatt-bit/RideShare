package db

import (
	"database/sql"
	"rideshare-go/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func InitMySQL(log *logrus.Logger, cfg config.DBConfig) (*sql.DB, error) {
	var err error
	db, err := sql.Open("mysql", cfg.DBConnectionString())
	if err != nil {
		return nil, err
	}

	// Connection test
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Info("connected to MySQL DB on ", cfg.Host, ":", cfg.Port)

	return db, nil
}
