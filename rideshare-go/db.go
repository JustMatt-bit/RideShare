package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func initMySQL(log *logrus.Logger, cfg DBConfig) (*sql.DB, error) {
	var err error
	db, err := sql.Open("mysql", cfg.dBConnectionString())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Info("connected to MySQL DB on ", cfg.Host, ":", cfg.Port)

	return db, nil
}

// ========================
// Ride Handlers
// ========================

func dbGetRideByID(id int) (ride, error) {
	row := Connection.QueryRow("SELECT * FROM ride WHERE id = ?", id)
	var r ride
	err := row.Scan(&r.ID, &r.OwnerID, &r.VehicleID, &r.StartDate, &r.StartCity, &r.StartAddress, &r.EndCity, &r.EndAddress, &r.CreatedAt)
	if err != nil {
		return ride{}, err
	}
	return r, nil
}

func dbGetRides() ([]ride, error) {
	rows, err := Connection.Query("SELECT * FROM ride")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []ride
	for rows.Next() {
		var r ride
		if err := rows.Scan(&r.ID, &r.OwnerID, &r.VehicleID, &r.StartDate, &r.StartCity, &r.StartAddress, &r.EndCity, &r.EndAddress, &r.CreatedAt); err != nil {
			return nil, err
		}
		rides = append(rides, r)
	}

	return rides, nil
}

func dbCreateRide(r ride) (int64, error) {
	result, err := Connection.Exec("INSERT INTO ride (owner_user_id, vehicle_id, start_date, start_city, start_address, end_city, end_address) VALUES (?, ?, ?, ?, ?, ?, ?)",
		r.OwnerID, r.VehicleID, r.StartDate, r.StartCity, r.StartAddress, r.EndCity, r.EndAddress)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func dbUpdateRide(id int, r ride) error {
	_, err := Connection.Exec("UPDATE ride SET owner_user_id = ?, vehicle_id = ?, start_date = ?, start_city = ?, start_address = ?, end_city = ?, end_address = ? WHERE id = ?",
		r.OwnerID, r.VehicleID, r.StartDate, r.StartCity, r.StartAddress, r.EndCity, r.EndAddress, id)
	return err
}

func dbDeleteRide(id int) error {
	_, err := Connection.Exec("DELETE FROM ride WHERE id = ?", id)
	return err
}

// ========================
// User Handlers
// ========================

func dbGetUsers() ([]user, error) {
	rows, err := Connection.Query("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Settings, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func dbGetUserByID(id int) (user, error) {
	row := Connection.QueryRow("SELECT * FROM user WHERE id = ?", id)
	var u user
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Settings, &u.CreatedAt)
	if err != nil {
		return user{}, err
	}
	return u, nil
}

func dbCreateUser(u user) (int64, error) {
	result, err := Connection.Exec("INSERT INTO user (name, email, password, settings) VALUES (?, ?, ?, ?)",
		u.Name, u.Email, u.Password, u.Settings)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func dbUpdateUser(id int, u user) error {
	_, err := Connection.Exec("UPDATE user SET name = ?, email = ?, password = ?, settings = ? WHERE id = ?",
		u.Name, u.Email, u.Password, u.Settings, id)
	return err
}

func dbDeleteUser(id int) error {
	_, err := Connection.Exec("DELETE FROM user WHERE id = ?", id)
	return err
}

func sqlErrorToHTTP(err error) (string, int) {
	if err == sql.ErrNoRows {
		return "not found", http.StatusNotFound
	}
	return "", http.StatusInternalServerError
}
