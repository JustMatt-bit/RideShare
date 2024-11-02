package db

import (
	"database/sql"
	"main/core"
)

func GetUserAuthByToken(db *sql.DB, token, authService string) (*core.UserAuthRecord, error) {
	row := db.QueryRow("SELECT * FROM auth WHERE token = ? AND auth_service = ?", token, authService)
	var ua core.UserAuthRecord
	if err := row.Scan(&ua.Token, &ua.Service, &ua.UserID); err != nil {
		return nil, err
	}
	return &ua, nil
}

func CreateUserAuth(db *sql.DB, ua core.UserAuthRecord) error {
	_, err := db.Exec("INSERT INTO auth (user_id, auth_service, token) VALUES (?, ?, ?)", ua.UserID, ua.Service, ua.Token)
	return err
}
