package db

import (
	"database/sql"
	"main/core"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func InitMySQL(log *logrus.Logger, cfg core.DBConfig) (*sql.DB, error) {
	var err error
	db, err := sql.Open("mysql", cfg.DBConnectionString())
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

// ========================
// User Handlers
// ========================

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

// ========================
// Car Handlers
// ========================

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
	result, err := db.Exec("INSERT INTO car (user_id, licence_plate, year, model_id) VALUES (?, ?, ?, ?)",
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
	_, err := db.Exec("UPDATE car SET user_id = ?, licence_platem = ?, year = ?, model_id = ? WHERE id = ?",
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

// ========================
// Car Model Handlers
// ========================

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

// ========================
// Car Make Handlers
// ========================

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

// ========================
// Car Category Handlers
// ========================

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

// ========================
// Ride Passenger Handlers
// ========================

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

// ========================
// Feedback Handlers
// ========================

func GetFeedbacks(db *sql.DB) ([]core.Feedback, error) {
	rows, err := db.Query("SELECT * FROM user_feedback")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []core.Feedback
	for rows.Next() {
		var f core.Feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func GetFeedbackByID(db *sql.DB, id int) (*core.Feedback, error) {
	row := db.QueryRow("SELECT * FROM user_feedback WHERE id = ?", id)
	var f core.Feedback
	if err := row.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
		return nil, err
	}
	return &f, nil
}

func GetFeedbacksByUserID(db *sql.DB, userID int) ([]core.Feedback, error) {
	rows, err := db.Query("SELECT * FROM user_feedback WHERE owner_user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []core.Feedback
	for rows.Next() {
		var f core.Feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func GetFeedbacksByRideID(db *sql.DB, rideID int) ([]core.Feedback, error) {
	rows, err := db.Query("SELECT * FROM user_feedback WHERE ride_id = ?", rideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []core.Feedback
	for rows.Next() {
		var f core.Feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func GetFeedbackByUserIDAndRideID(db *sql.DB, userID, rideID int) ([]core.Feedback, error) {
	rows, err := db.Query("SELECT uf.id, uf.owner_user_id, uf.ride_id, uf.score, uf.message, uf.created_at FROM user_feedback uf LEFT JOIN ride r ON r.id = uf.ride_id WHERE r.owner_user_id = ? AND uf.ride_id = ?", userID, rideID)
	if err != nil {
		return nil, err
	}

	var feedbacks []core.Feedback
	for rows.Next() {
		var f core.Feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.RideID, &f.Score, &f.Message, &f.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}

	return feedbacks, nil
}

func CreateFeedback(db *sql.DB, f core.Feedback) (int64, error) {
	result, err := db.Exec("INSERT INTO user_feedback (owner_user_id, ride_id, score, message) VALUES (?, ?, ?, ?)", f.UserID, f.RideID, f.Score, f.Message)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateFeedback(db *sql.DB, id int, f core.Feedback) error {
	_, err := db.Exec("UPDATE user_feedback SET owner_user_id = ?, ride_id = ?, score = ?, message = ? WHERE id = ?",
		f.UserID, f.RideID, f.Score, f.Message, id)
	return err
}

func DeleteFeedback(db *sql.DB, feedbackID int) error {
	_, err := db.Exec("DELETE FROM user_feedback WHERE id = ?", feedbackID)
	return err
}

// ========================
// Auth Functions
// ========================

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

func SqlErrorToHTTP(err error) (string, int) {
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
