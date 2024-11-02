package db

import (
	"database/sql"
	"main/core"
)

func GetCarMakes(db *sql.DB) ([]core.CarMake, error) {
	rows, err := db.Query("SELECT * FROM car_make")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var makes []core.CarMake
	for rows.Next() {
		var m core.CarMake
		if err := rows.Scan(&m.ID, &m.Name); err != nil {
			return nil, err
		}
		makes = append(makes, m)
	}

	return makes, nil
}

func GetCarMakeByID(db *sql.DB, id int) (*core.CarMake, error) {
	row := db.QueryRow("SELECT * FROM car_make WHERE id = ?", id)
	var m core.CarMake
	if err := row.Scan(&m.ID, &m.Name); err != nil {
		return nil, err
	}
	return &m, nil
}

func CreateCarMake(db *sql.DB, m core.CarMake) (int64, error) {
	result, err := db.Exec("INSERT INTO car_make (name) VALUES (?)", m.Name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateCarMake(db *sql.DB, id int, m core.CarMake) error {
	_, err := db.Exec("UPDATE car_make SET name = ? WHERE id = ?", m.Name, id)
	return err
}

func DeleteCarMake(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM car_make WHERE id = ?", id)
	return err
}
