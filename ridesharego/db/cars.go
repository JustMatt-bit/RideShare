package db

import (
	"database/sql"
	"main/core"
)

func GetCars(db *sql.DB) ([]core.Car, error) {
	rows, err := db.Query("SELECT * FROM car")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []core.Car
	for rows.Next() {
		var c core.Car
		if err := rows.Scan(&c.ID, &c.LicensePlate, &c.UserID, &c.ModelID, &c.Year); err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}

	return cars, nil
}

func GetCarByID(db *sql.DB, id int) (*core.Car, error) {
	row := db.QueryRow("SELECT * FROM car WHERE id = ?", id)
	var c core.Car
	if err := row.Scan(&c.ID, &c.LicensePlate, &c.UserID, &c.ModelID, &c.Year); err != nil {
		return nil, err
	}
	return &c, nil
}

func CreateCar(db *sql.DB, c core.Car) (int64, error) {
	result, err := db.Exec("INSERT INTO car (user_id, license_plate, year, model_id) VALUES (?, ?, ?, ?)",
		c.UserID, c.LicensePlate, c.Year, c.ModelID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateCar(db *sql.DB, id int, c core.Car) error {
	_, err := db.Exec("UPDATE car SET user_id = ?, license_plate = ?, year = ?, model_id = ? WHERE id = ?",
		c.UserID, c.LicensePlate, c.Year, c.ModelID, id)
	return err
}

func DeleteCar(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM car WHERE id = ?", id)
	return err
}

func GetCarsByUserID(db *sql.DB, userID int) ([]core.Car, error) {
	rows, err := db.Query("SELECT * FROM car WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []core.Car
	for rows.Next() {
		var c core.Car
		if err := rows.Scan(&c.ID, &c.LicensePlate, &c.UserID, &c.ModelID, &c.Year); err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}

	return cars, nil
}
