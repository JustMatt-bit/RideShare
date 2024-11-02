package db

import (
	"database/sql"
	"main/core"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func InitMySQL(log *logrus.Entry, cfg core.DBConfig) (*sql.DB, error) {
	var err error
	db, err := sql.Open("mysql", cfg.DBConnectionString())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Info("connected to MySQL DB on ", cfg.Host, ":", cfg.Port)

	return db, nil
}

func SqlErrorToHTTP(err error) (string, int) {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 {
			return "duplicate entry", http.StatusConflict
		}
	}
	if err == sql.ErrNoRows {
		return "not found", http.StatusNotFound
	}
	return "internal sever error", http.StatusInternalServerError
}
