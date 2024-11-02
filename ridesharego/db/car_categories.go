package db

import (
	"database/sql"
	"main/core"
)

func GetCarCategories(db *sql.DB) ([]core.CarCategory, error) {
	rows, err := db.Query("SELECT * FROM car_category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []core.CarCategory
	for rows.Next() {
		var c core.CarCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.PassengerCount); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func GetCarCategoryByID(db *sql.DB, id int) (*core.CarCategory, error) {
	row := db.QueryRow("SELECT * FROM car_category WHERE id = ?", id)
	var c core.CarCategory
	if err := row.Scan(&c.ID, &c.Name, &c.PassengerCount); err != nil {
		return nil, err
	}
	return &c, nil
}

func CreateCarCategory(db *sql.DB, c core.CarCategory) (int64, error) {
	result, err := db.Exec("INSERT INTO car_category (name, passenger_count) VALUES (?, ?)", c.Name, c.PassengerCount)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateCarCategory(db *sql.DB, id int, c core.CarCategory) error {
	_, err := db.Exec("UPDATE car_category SET name = ?, passenger_count = ? WHERE id = ?", c.Name, c.PassengerCount, id)
	return err
}

func DeleteCarCategory(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM car_category WHERE id = ?", id)
	return err
}
