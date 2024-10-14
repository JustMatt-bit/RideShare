package main

import (
	"database/sql"
	"net/http"

	"github.com/go-sql-driver/mysql"
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

func dbGetRidesByUserID(userID int) ([]ride, error) {
	rows, err := Connection.Query("SELECT * FROM ride WHERE owner_user_id = ?", userID)
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

// ========================
// Car Handlers
// ========================

func dbGetCars() ([]car, error) {
	rows, err := Connection.Query("SELECT * FROM car")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []car
	for rows.Next() {
		var c car
		if err := rows.Scan(&c.ID, &c.LicensePlate, &c.UserID, &c.ModelID, &c.Year); err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}

	return cars, nil
}

func dbGetCarByID(id int) (car, error) {
	row := Connection.QueryRow("SELECT * FROM car WHERE id = ?", id)
	var c car
	if err := row.Scan(&c.ID, &c.LicensePlate, &c.UserID, &c.ModelID, &c.Year); err != nil {
		return car{}, err
	}
	return c, nil
}

func dbCreateCar(c car) (int64, error) {
	result, err := Connection.Exec("INSERT INTO car (user_id, licence_plate, year, model_id) VALUES (?, ?, ?, ?)",
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

func dbUpdateCar(id int, c car) error {
	_, err := Connection.Exec("UPDATE car SET user_id = ?, licence_platem = ?, year = ?, model_id = ? WHERE id = ?",
		c.UserID, c.LicensePlate, c.Year, c.ModelID, id)
	return err
}

func dbDeleteCar(id int) error {
	_, err := Connection.Exec("DELETE FROM car WHERE id = ?", id)
	return err
}

// ========================
// Car Model Handlers
// ========================

func dbGetCarModels() ([]carModel, error) {
	rows, err := Connection.Query("SELECT * FROM car_model")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []carModel
	for rows.Next() {
		var m carModel
		if err := rows.Scan(&m.ID, &m.CategoryID, &m.MakeID, &m.Name); err != nil {
			return nil, err
		}
		models = append(models, m)
	}

	return models, nil
}

func dbGetCarModelByID(id int) (carModel, error) {
	row := Connection.QueryRow("SELECT * FROM car_model WHERE id = ?", id)
	var m carModel
	if err := row.Scan(&m.ID, &m.CategoryID, &m.MakeID, &m.Name); err != nil {
		return carModel{}, err
	}
	return m, nil
}

func dbCreateCarModel(m carModel) (int64, error) {
	result, err := Connection.Exec("INSERT INTO car_model (category_id, make_id, name) VALUES (?, ?, ?)", m.CategoryID, m.MakeID, m.Name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func dbUpdateCarModel(id int, m carModel) error {
	_, err := Connection.Exec("UPDATE car_model SET category_id = ?, make_id = ?, name = ? WHERE id = ?", m.CategoryID, m.MakeID, m.Name, id)
	return err
}

func dbDeleteCarModel(id int) error {
	_, err := Connection.Exec("DELETE FROM car_model WHERE id = ?", id)
	return err
}

// ========================
// Car Make Handlers
// ========================

func dbGetCarMakes() ([]carMake, error) {
	rows, err := Connection.Query("SELECT * FROM car_make")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var makes []carMake
	for rows.Next() {
		var m carMake
		if err := rows.Scan(&m.ID, &m.Name); err != nil {
			return nil, err
		}
		makes = append(makes, m)
	}

	return makes, nil
}

func dbGetCarMakeByID(id int) (carMake, error) {
	row := Connection.QueryRow("SELECT * FROM car_make WHERE id = ?", id)
	var m carMake
	if err := row.Scan(&m.ID, &m.Name); err != nil {
		return carMake{}, err
	}
	return m, nil
}

func dbCreateCarMake(m carMake) (int64, error) {
	result, err := Connection.Exec("INSERT INTO car_make (name) VALUES (?)", m.Name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func dbUpdateCarMake(id int, m carMake) error {
	_, err := Connection.Exec("UPDATE car_make SET name = ? WHERE id = ?", m.Name, id)
	return err
}

func dbDeleteCarMake(id int) error {
	_, err := Connection.Exec("DELETE FROM car_make WHERE id = ?", id)
	return err
}

// ========================
// Car Category Handlers
// ========================

func dbGetCarCategories() ([]carCategory, error) {
	rows, err := Connection.Query("SELECT * FROM car_category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []carCategory
	for rows.Next() {
		var c carCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.PassengerCount); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func dbGetCarCategoryByID(id int) (carCategory, error) {
	row := Connection.QueryRow("SELECT * FROM car_category WHERE id = ?", id)
	var c carCategory
	if err := row.Scan(&c.ID, &c.Name, &c.PassengerCount); err != nil {
		return carCategory{}, err
	}
	return c, nil
}

func dbCreateCarCategory(c carCategory) (int64, error) {
	result, err := Connection.Exec("INSERT INTO car_category (name, passenger_count) VALUES (?, ?)", c.Name, c.PassengerCount)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func dbUpdateCarCategory(id int, c carCategory) error {
	_, err := Connection.Exec("UPDATE car_category SET name = ?, passenger_count = ? WHERE id = ?", c.Name, c.PassengerCount, id)
	return err
}

func dbDeleteCarCategory(id int) error {
	_, err := Connection.Exec("DELETE FROM car_category WHERE id = ?", id)
	return err
}

// ========================
// Ride Passenger Handlers
// ========================

func dbGetPassengersByRideID(rideID int) ([]passenger, error) {
	rows, err := Connection.Query("SELECT * FROM ride_passenger WHERE ride_id = ?", rideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passengers []passenger
	for rows.Next() {
		var rp passenger
		if err := rows.Scan(&rp.RideID, &rp.PassengerID, &rp.CreatedAt); err != nil {
			return nil, err
		}
		passengers = append(passengers, rp)
	}

	return passengers, nil
}

func dbGetPassengerByRideIDAndUserID(rideID, userID int) (passenger, error) {
	row := Connection.QueryRow("SELECT * FROM ride_passenger WHERE ride_id = ? AND passenger_id = ?", rideID, userID)
	var rp passenger
	if err := row.Scan(&rp.RideID, &rp.PassengerID, &rp.CreatedAt); err != nil {
		return passenger{}, err
	}
	return rp, nil
}

func dbCreatePassenger(rp passenger) error {
	_, err := Connection.Exec("INSERT INTO ride_passenger (ride_id, passenger_id) VALUES (?, ?)", rp.RideID, rp.PassengerID)
	if err != nil {
		return err
	}

	return nil
}

func dbDeletePassenger(rideID, userID int) error {
	_, err := Connection.Exec("DELETE FROM ride_passenger WHERE ride_id = ? AND passenger_id = ?", rideID, userID)
	return err
}

// ========================
// Feedback Handlers
// ========================

func dbGetFeedbacks() ([]feedback, error) {
	rows, err := Connection.Query("SELECT * FROM user_feedback")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []feedback
	for rows.Next() {
		var f feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func dbGetFeedbackByID(id int) (feedback, error) {
	row := Connection.QueryRow("SELECT * FROM user_feedback WHERE id = ?", id)
	var f feedback
	if err := row.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
		return feedback{}, err
	}
	return f, nil
}

func dbGetFeedbacksByUserID(userID int) ([]feedback, error) {
	rows, err := Connection.Query("SELECT * FROM user_feedback WHERE owner_user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []feedback
	for rows.Next() {
		var f feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func dbGetFeedbacksByRideID(rideID int) ([]feedback, error) {
	rows, err := Connection.Query("SELECT * FROM user_feedback WHERE ride_id = ?", rideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []feedback
	for rows.Next() {
		var f feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func dbGetFeedbackByUserIDAndRideID(userID, rideID int) ([]feedback, error) {
	rows, err := Connection.Query("SELECT uf.id, uf.owner_user_id, uf.ride_id, uf.score, uf.message, uf.created_at FROM user_feedback uf LEFT JOIN ride r ON r.id = uf.ride_id WHERE r.owner_user_id = ? AND uf.ride_id = ?", userID, rideID)
	if err != nil {
		return nil, err
	}

	var feedbacks []feedback
	for rows.Next() {
		var f feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func dbCreateFeedback(f feedback) (int64, error) {
	result, err := Connection.Exec("INSERT INTO user_feedback (owner_user_id, ride_id, score, message) VALUES (?, ?, ?, ?)", f.UserID, f.RideID, f.Score, f.Message)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func dbUpdateFeedback(id int, f feedback) error {
	_, err := Connection.Exec("UPDATE user_feedback SET owner_user_id = ?, ride_id = ?, score = ?, message = ? WHERE id = ?",
		f.UserID, f.RideID, f.Score, f.Message, id)
	return err
}

func dbDeleteFeedback(feedbackID int) error {
	_, err := Connection.Exec("DELETE FROM user_feedback WHERE id = ?", feedbackID)
	return err
}

func sqlErrorToHTTP(err error) (string, int) {
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
