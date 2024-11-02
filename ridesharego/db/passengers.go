package db

import (
	"database/sql"
	"main/core"
)

func GetPassengersByRideID(db *sql.DB, rideID int) ([]core.Passenger, error) {
	rows, err := db.Query("SELECT * FROM ride_passenger WHERE ride_id = ?", rideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passengers []core.Passenger
	for rows.Next() {
		var rp core.Passenger
		if err := rows.Scan(&rp.RideID, &rp.PassengerID, &rp.CreatedAt); err != nil {
			return nil, err
		}
		passengers = append(passengers, rp)
	}

	return passengers, nil
}

func GetPassengerByRideIDAndUserID(db *sql.DB, rideID, userID int) (*core.Passenger, error) {
	row := db.QueryRow("SELECT * FROM ride_passenger WHERE ride_id = ? AND passenger_id = ?", rideID, userID)
	var rp core.Passenger
	if err := row.Scan(&rp.RideID, &rp.PassengerID, &rp.CreatedAt); err != nil {
		return nil, err
	}
	return &rp, nil
}

func CreatePassenger(db *sql.DB, rp core.Passenger) error {
	_, err := db.Exec("INSERT INTO ride_passenger (ride_id, passenger_id) VALUES (?, ?)", rp.RideID, rp.PassengerID)
	if err != nil {
		return err
	}

	return nil
}

func DeletePassenger(db *sql.DB, rideID, userID int) error {
	_, err := db.Exec("DELETE FROM ride_passenger WHERE ride_id = ? AND passenger_id = ?", rideID, userID)
	return err
}
