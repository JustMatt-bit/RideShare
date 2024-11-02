package db

import (
	"database/sql"
	"main/core"
)

func GetUsers(db *sql.DB) ([]core.User, error) {
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []core.User
	for rows.Next() {
		var u core.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Settings, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func GetUserByID(db *sql.DB, id int64) (*core.User, error) {
	row := db.QueryRow("SELECT * FROM user WHERE id = ?", id)
	var u core.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Settings, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUser(db *sql.DB, u core.User) (int64, error) {
	settings := u.Settings
	if len(settings) == 0 {
		settings = []byte("{}")
	}

	result, err := db.Exec("INSERT INTO user (name, email, password, settings) VALUES (?, ?, ?, ?)",
		u.Name, u.Email, u.Password, settings)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateUser(db *sql.DB, id int64, u core.User) error {
	_, err := db.Exec("UPDATE user SET name = ?, email = ?, password = ?, role = ?, settings = ? WHERE id = ?",
		u.Name, u.Email, u.Password, u.Role, u.Settings, id)
	return err
}

func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM user WHERE id = ?", id)
	return err
}
