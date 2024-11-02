package db

import (
	"database/sql"
	"main/core"
)

func GetRideByID(db *sql.DB, id int) (*core.Ride, error) {
	row := db.QueryRow("SELECT * FROM ride WHERE id = ?", id)
	var r core.Ride
	err := row.Scan(&r.ID, &r.OwnerID, &r.VehicleID, &r.StartDate, &r.StartCity, &r.StartAddress, &r.EndCity, &r.EndAddress, &r.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func GetRides(db *sql.DB) ([]core.Ride, error) {
	rows, err := db.Query("SELECT * FROM ride")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []core.Ride
	for rows.Next() {
		var r core.Ride
		if err := rows.Scan(&r.ID, &r.OwnerID, &r.VehicleID, &r.StartDate, &r.StartCity, &r.StartAddress, &r.EndCity, &r.EndAddress, &r.CreatedAt); err != nil {
			return nil, err
		}
		rides = append(rides, r)
	}

	return rides, nil
}

func CreateRide(db *sql.DB, r core.Ride) (int64, error) {
	result, err := db.Exec("INSERT INTO ride (owner_user_id, vehicle_id, start_date, start_city, start_address, end_city, end_address) VALUES (?, ?, ?, ?, ?, ?, ?)",
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

func UpdateRide(db *sql.DB, id int, r core.Ride) error {
	_, err := db.Exec("UPDATE ride SET owner_user_id = ?, vehicle_id = ?, start_date = ?, start_city = ?, start_address = ?, end_city = ?, end_address = ? WHERE id = ?",
		r.OwnerID, r.VehicleID, r.StartDate, r.StartCity, r.StartAddress, r.EndCity, r.EndAddress, id)
	return err
}

func DeleteRide(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM ride WHERE id = ?", id)
	return err
}

func GetRidesByUserID(db *sql.DB, userID int) ([]core.Ride, error) {
	rows, err := db.Query("SELECT * FROM ride WHERE owner_user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []core.Ride
	for rows.Next() {
		var r core.Ride
		if err := rows.Scan(&r.ID, &r.OwnerID, &r.VehicleID, &r.StartDate, &r.StartCity, &r.StartAddress, &r.EndCity, &r.EndAddress, &r.CreatedAt); err != nil {
			return nil, err
		}
		rides = append(rides, r)
	}

	return rides, nil
}
