package db

import (
	"database/sql"
	"main/core"
)

func GetCarModels(db *sql.DB) ([]core.CarModel, error) {
	rows, err := db.Query("SELECT * FROM car_model")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []core.CarModel
	for rows.Next() {
		var m core.CarModel
		if err := rows.Scan(&m.ID, &m.CategoryID, &m.MakeID, &m.Name); err != nil {
			return nil, err
		}
		models = append(models, m)
	}

	return models, nil
}

func GetCarModelByID(db *sql.DB, id int) (*core.CarModel, error) {
	row := db.QueryRow("SELECT * FROM car_model WHERE id = ?", id)
	var m core.CarModel
	if err := row.Scan(&m.ID, &m.CategoryID, &m.MakeID, &m.Name); err != nil {
		return nil, err
	}
	return &m, nil
}

func CreateCarModel(db *sql.DB, m core.CarModel) (int64, error) {
	result, err := db.Exec("INSERT INTO car_model (category_id, make_id, name) VALUES (?, ?, ?)", m.CategoryID, m.MakeID, m.Name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateCarModel(db *sql.DB, id int, m core.CarModel) error {
	_, err := db.Exec("UPDATE car_model SET category_id = ?, make_id = ?, name = ? WHERE id = ?", m.CategoryID, m.MakeID, m.Name, id)
	return err
}

func DeleteCarModel(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM car_model WHERE id = ?", id)
	return err
}
