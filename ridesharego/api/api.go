package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"main/core"
	"main/db"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const responseTypeXML = "application/xml"

func CreateRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.Use(loggerMiddleware)

	dbMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), core.CtxDB, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	r.Use(dbMiddleware)

	// Google OAuth2 endpoint
	api := r.PathPrefix("/api/v1").Subrouter()

	withAdmin := func(handler func(http.ResponseWriter, *http.Request, *logrus.Logger, *sql.DB)) http.HandlerFunc {
		return authMiddleware(core.RoleAdmin, withMiddleware(handler))
	}

	withUser := func(handler func(http.ResponseWriter, *http.Request, *logrus.Logger, *sql.DB)) http.HandlerFunc {
		return authMiddleware(core.RoleUser, withMiddleware(handler))
	}

	withGuest := func(handler func(http.ResponseWriter, *http.Request, *logrus.Logger, *sql.DB)) http.HandlerFunc {
		return withMiddleware(handler)
	}

	// Ride endpoints
	api.HandleFunc("/rides", withGuest(getRides)).Methods("GET")
	api.HandleFunc("/ride/{ride_id}", withGuest(getRide)).Methods("GET")
	api.HandleFunc("/ride", withUser(createRide)).Methods("POST")
	api.HandleFunc("/ride/{ride_id}", withUser(updateRide)).Methods("PUT")
	api.HandleFunc("/ride/{ride_id}", withUser(deleteRide)).Methods("DELETE")
	api.HandleFunc("/user/{user_id}/rides", withUser(getUserRides)).Methods("GET")

	// User endpoints
	api.HandleFunc("/users", withAdmin(getUsers)).Methods("GET")
	api.HandleFunc("/user/{user_id}", withAdmin(getUser)).Methods("GET")
	api.HandleFunc("/user", withAdmin(createUser)).Methods("POST")
	api.HandleFunc("/user/{user_id}", withUser(updateUser)).Methods("PUT")
	api.HandleFunc("/user/{user_id}", withUser(deleteUser)).Methods("DELETE")

	// Car endpoints
	api.HandleFunc("/cars", withAdmin(getCars)).Methods("GET")
	api.HandleFunc("/car/{car_id}", withAdmin(getCar)).Methods("GET")
	api.HandleFunc("/car", withUser(createCar)).Methods("POST")
	api.HandleFunc("/car/{car_id}", withUser(updateCar)).Methods("PUT")
	api.HandleFunc("/car/{car_id}", withUser(deleteCar)).Methods("DELETE")
	api.HandleFunc("/user/{user_id}/cars", withUser(getUserCars)).Methods("GET")

	// Car model endpoints
	api.HandleFunc("/car_models", withAdmin(getCarModels)).Methods("GET")
	api.HandleFunc("/car_model/{model_id}", withAdmin(getCarModel)).Methods("GET")
	api.HandleFunc("/car_model", withAdmin(createCarModel)).Methods("POST")
	api.HandleFunc("/car_model/{model_id}", withAdmin(updateCarModel)).Methods("PUT")
	api.HandleFunc("/car_model/{model_id}", withAdmin(deleteCarModel)).Methods("DELETE")

	// Car make endpoints
	api.HandleFunc("/car_makes", withAdmin(getCarMakes)).Methods("GET")
	api.HandleFunc("/car_make/{make_id}", withAdmin(getCarMake)).Methods("GET")
	api.HandleFunc("/car_make", withAdmin(createCarMake)).Methods("POST")
	api.HandleFunc("/car_make/{make_id}", withAdmin(updateCarMake)).Methods("PUT")
	api.HandleFunc("/car_make/{make_id}", withAdmin(deleteCarMake)).Methods("DELETE")

	// Car category endpoints
	api.HandleFunc("/car_categories", withAdmin(getCarCategories)).Methods("GET")
	api.HandleFunc("/car_category/{category_id}", withAdmin(getCarCategory)).Methods("GET")
	api.HandleFunc("/car_category", withAdmin(createCarCategory)).Methods("POST")
	api.HandleFunc("/car_category/{category_id}", withAdmin(updateCarCategory)).Methods("PUT")
	api.HandleFunc("/car_category/{category_id}", withAdmin(deleteCarCategory)).Methods("DELETE")

	// Ride passenger endpoints
	api.HandleFunc("/ride/{ride_id}/passengers", withUser(getRidePassengers)).Methods("GET")
	api.HandleFunc("/ride/{ride_id}/passenger/{user_id}", withUser(createRidePassenger)).Methods("POST")
	api.HandleFunc("/ride/{ride_id}/passenger/{user_id}", withUser(deleteRidePassenger)).Methods("DELETE")

	// Feedback endpoints
	api.HandleFunc("/feedback", withAdmin(getFeedbacks)).Methods("GET")
	api.HandleFunc("/feedback/{feedback_id}", withAdmin(getFeedback)).Methods("GET")
	api.HandleFunc("/feedback", withUser(createFeedback)).Methods("POST")
	api.HandleFunc("/feedback/{feedback_id}", withUser(updateFeedback)).Methods("PUT")
	api.HandleFunc("/feedback/{feedback_id}", withUser(deleteFeedback)).Methods("DELETE")
	api.HandleFunc("/ride/{ride_id}/feedback", withGuest(getRideFeedback)).Methods("GET")
	api.HandleFunc("/user/{user_id}/feedback", withUser(getUserFeedback)).Methods("GET")
	api.HandleFunc("/user/{user_id}/ride/{ride_id}/feedback", withUser(getUserRideFeedback)).Methods("GET")

	return r
}

// =====================
//   Ride Handlers
// =====================

func getRides(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	rides, err := db.GetRides(d)
	if err != nil {
		log.WithError(err).Error("getting rides")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, rides)
}

func getRide(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["ride_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ride, err := db.GetRideByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, ride)
}

func createRide(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	var ride core.Ride
	if err := json.NewDecoder(r.Body).Decode(&ride); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if ride.OwnerID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	_, err := db.GetUserByID(d, int64(ride.OwnerID))
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting user: %s", error), status)
		return
	}

	car, err := db.GetCarByID(d, ride.VehicleID)
	if err != nil {
		log.WithError(err).Error("getting car")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting car: %s", error), status)
		return
	}

	if err := ride.Validate(car); err != nil {
		http.Error(w, fmt.Sprintf("validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := db.CreateRide(d, ride)
	if err != nil {
		log.WithError(err).Error("creating ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	ride.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, ride)
}

func updateRide(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["ride_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ride core.Ride
	if err := json.NewDecoder(r.Body).Decode(&ride); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if ride.OwnerID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	_, err = db.GetUserByID(d, int64(ride.OwnerID))
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting user: %s", error), status)
		return
	}

	car, err := db.GetCarByID(d, ride.VehicleID)
	if err != nil {
		log.WithError(err).Error("getting car")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting car: %s", error), status)
		return
	}

	if err := ride.Validate(car); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := db.GetRideByID(d, idInt); err != nil {
		log.WithError(err).Error("getting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting ride: %s", error), status)
		return
	}

	if err := db.UpdateRide(d, idInt, ride); err != nil {
		log.WithError(err).Error("updating ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteRide(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["ride_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ride, err := db.GetRideByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting ride: %s", error), status)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if userAuth.UserID != ride.OwnerID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	if err := db.DeleteRide(d, idInt); err != nil {
		log.WithError(err).Error("deleting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   User Handlers
// =====================

func getUsers(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	users, err := db.GetUsers(d)
	if err != nil {
		log.WithError(err).Error("getting users")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, users)
}

func getUser(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["user_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(d, int64(idInt))
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, user)
}

func createUser(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	var user core.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := db.CreateUser(d, user)
	if err != nil {
		log.WithError(err).Error("creating user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	user.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, user)
}

func updateUser(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["user_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var userUpdate core.User
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := userUpdate.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	existingUser, err := db.GetUserByID(d, int64(idInt))
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting user: %s", error), status)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if existingUser.Role != userUpdate.Role && userAuth.UserID != userUpdate.ID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "invalid permissions", http.StatusForbidden)
		return
	}

	if err := db.UpdateUser(d, int64(idInt), userUpdate); err != nil {
		log.WithError(err).Error("updating user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteUser(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["user_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	existingUser, err := db.GetUserByID(d, int64(idInt))
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting user: %s", error), status)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if existingUser.ID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	if err := db.DeleteUser(d, idInt); err != nil {
		log.WithError(err).Error("deleting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getUserRides(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["user_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if idInt != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	rides, err := db.GetRidesByUserID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting user rides")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, rides)
}

func getUserCars(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["user_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if idInt != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	cars, err := db.GetCarsByUserID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting user rides")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, cars)
}

// =====================
//   Car Handlers
// =====================

func getCars(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	cars, err := db.GetCars(d)
	if err != nil {
		log.WithError(err).Error("getting cars")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, cars)
}

func getCar(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["car_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	car, err := db.GetCarByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting car")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, car)
}

func createCar(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	var car core.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if car.UserID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	_, err := db.GetUserByID(d, int64(car.UserID))
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting user: %s", error), status)
		return
	}

	_, err = db.GetCarModelByID(d, car.ModelID)
	if err != nil {
		log.WithError(err).Error("getting car model")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting car model: %s", error), status)
		return
	}

	cars, err := db.GetCarsByUserID(d, car.UserID)
	if err != nil {
		log.WithError(err).Error("getting cars")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting cars: %s", error), status)
		return
	}

	if err := car.Validate(cars); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := db.CreateCar(d, car)
	if err != nil {
		log.WithError(err).Error("creating car")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	car.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, car)
}

func updateCar(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["car_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var car core.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if car.UserID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	_, err = db.GetUserByID(d, int64(car.UserID))
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting user: %s", error), status)
		return
	}

	_, err = db.GetCarModelByID(d, car.ModelID)
	if err != nil {
		log.WithError(err).Error("getting car model")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting car model: %s", error), status)
		return
	}

	cars, err := db.GetCarsByUserID(d, car.UserID)
	if err != nil {
		log.WithError(err).Error("getting cars")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting cars: %s", error), status)
		return
	}

	if err := car.Validate(cars); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := db.GetCarByID(d, idInt); err != nil {
		log.WithError(err).Error("getting car")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car: %s", error), status)
		return
	}

	if err := db.UpdateCar(d, idInt, car); err != nil {
		log.WithError(err).Error("updating car")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteCar(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["car_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if idInt != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	if _, err := db.GetCarByID(d, idInt); err != nil {
		log.WithError(err).Error("getting car")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car: %s", error), status)
		return
	}

	if err := db.DeleteCar(d, idInt); err != nil {
		log.WithError(err).Error("deleting car")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   Car Model Handlers
// =====================

func getCarModels(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	models, err := db.GetCarModels(d)
	if err != nil {
		log.WithError(err).Error("getting car models")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, models)
}

func getCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["model_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	model, err := db.GetCarModelByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting car model")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, model)
}

func createCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	var model core.CarModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.GetCarCategoryByID(d, model.CategoryID)
	if err != nil {
		log.WithError(err).Error("getting car category")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting car category: %s", error), status)
		return
	}

	_, err = db.GetCarMakeByID(d, model.MakeID)
	if err != nil {
		log.WithError(err).Error("getting car make")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting car make: %s", error), status)
		return
	}

	if err := model.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := db.CreateCarModel(d, model)
	if err != nil {
		log.WithError(err).Error("creating car model")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	model.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, model)
}

func updateCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["model_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var model core.CarModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	_, err = db.GetCarCategoryByID(d, model.CategoryID)
	if err != nil {
		log.WithError(err).Error("getting car category")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting car category: %s", error), status)
		return
	}

	_, err = db.GetCarMakeByID(d, model.MakeID)
	if err != nil {
		log.WithError(err).Error("getting car make")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting car make: %s", error), status)
		return
	}

	if err := model.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := db.GetCarModelByID(d, idInt); err != nil {
		log.WithError(err).Error("getting car model")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car model: %s", error), status)
		return
	}

	if err := db.UpdateCarModel(d, idInt, model); err != nil {
		log.WithError(err).Error("updating car model")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["model_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if _, err := db.GetCarModelByID(d, idInt); err != nil {
		log.WithError(err).Error("getting car model")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car model: %s", error), status)
		return
	}

	if err := db.DeleteCarModel(d, idInt); err != nil {
		log.WithError(err).Error("deleting car model")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   Car Make Handlers
// =====================

func getCarMakes(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	makes, err := db.GetCarMakes(d)
	if err != nil {
		log.WithError(err).Error("getting car makes")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, makes)
}

func getCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["make_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	make, err := db.GetCarMakeByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting car make")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, make)
}

func createCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	var make core.CarMake
	if err := json.NewDecoder(r.Body).Decode(&make); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := make.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := db.CreateCarMake(d, make)
	if err != nil {
		log.WithError(err).Error("creating car make")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	make.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, make)
}

func updateCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["make_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var make core.CarMake
	if err := json.NewDecoder(r.Body).Decode(&make); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := make.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := db.GetCarMakeByID(d, idInt); err != nil {
		log.WithError(err).Error("getting car make")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car make: %s", error), status)
		return
	}

	if err := db.UpdateCarMake(d, idInt, make); err != nil {
		log.WithError(err).Error("updating car make")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["make_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if _, err := db.GetCarMakeByID(d, idInt); err != nil {
		log.WithError(err).Error("getting car make")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car make: %s", error), status)
		return
	}

	if err := db.DeleteCarMake(d, idInt); err != nil {
		log.WithError(err).Error("deleting car make")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   Car Category Handlers
// =====================

func getCarCategories(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	categories, err := db.GetCarCategories(d)
	if err != nil {
		log.WithError(err).Error("getting car categories")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, categories)
}

func getCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["category_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	category, err := db.GetCarCategoryByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting car category")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, category)
}

func createCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	var category core.CarCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := category.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := db.CreateCarCategory(d, category)
	if err != nil {
		log.WithError(err).Error("creating car category")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	category.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, category)
}

func updateCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["category_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var category core.CarCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := category.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := db.GetCarCategoryByID(d, idInt); err != nil {
		log.WithError(err).Error("getting car category")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car category: %s", error), status)
		return
	}

	if err := db.UpdateCarCategory(d, idInt, category); err != nil {
		log.WithError(err).Error("updating car category")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["category_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if _, err := db.GetCarCategoryByID(d, idInt); err != nil {
		log.WithError(err).Error("getting car category")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car category: %s", error), status)
		return
	}

	if err := db.DeleteCarCategory(d, idInt); err != nil {
		log.WithError(err).Error("deleting car category")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   Ride Passenger Handlers
// =====================

func getRidePassengers(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]

	if rideID == "" {
		http.Error(w, "missing ride_id", http.StatusBadRequest)
		return
	}

	rideIDInt, err := strconv.Atoi(rideID)
	if err != nil {
		http.Error(w, "invalid ride_id", http.StatusBadRequest)
		return
	}

	ride, err := db.GetRideByID(d, rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting ride: %s", error), status)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if ride.OwnerID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	passengers, err := db.GetPassengersByRideID(d, rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting ride passengers")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	passengerFound := false
	for _, passenger := range passengers {
		if passenger.PassengerID == userAuth.UserID {
			passengerFound = true
			break
		}
	}

	if !passengerFound {
		http.Error(w, "user is not a passenger", http.StatusForbidden)
		return
	}

	respond(w, r, passengers)
}

func createRidePassenger(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]
	userID := vars["user_id"]

	if rideID == "" {
		http.Error(w, "missing ride_id", http.StatusBadRequest)
		return
	}

	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	rideIDInt, err := strconv.Atoi(rideID)
	if err != nil {
		http.Error(w, "invalid ride_id", http.StatusBadRequest)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if userIDInt != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "user is not the passenger", http.StatusForbidden)
		return
	}

	passenger := core.Passenger{
		RideID:      rideIDInt,
		PassengerID: userIDInt,
	}

	_, err = db.GetUserByID(d, int64(userIDInt))
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting user: %s", error), status)
		return
	}

	ride, err := db.GetRideByID(d, rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting ride: %s", error), status)
		return
	}

	passengers, err := db.GetPassengersByRideID(d, rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting ride passengers")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting ride passengers: %s", error), status)
		return
	}

	car, err := db.GetCarByID(d, ride.VehicleID)
	if err != nil {
		log.WithError(err).Error("getting car")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car: %s", error), status)
		return
	}

	model, err := db.GetCarModelByID(d, car.ModelID)
	if err != nil {
		log.WithError(err).Error("getting car model")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car model: %s", error), status)
		return
	}

	category, err := db.GetCarCategoryByID(d, model.CategoryID)
	if err != nil {
		log.WithError(err).Error("getting car category")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car category: %s", error), status)
		return
	}

	if err := passenger.Validate(ride, passengers, category.PassengerCount); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = db.CreatePassenger(d, passenger)
	if err != nil {
		log.WithError(err).Error("creating ride passenger")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respond(w, r, passenger)
}

func deleteRidePassenger(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]
	userID := vars["user_id"]

	if rideID == "" {
		http.Error(w, "missing ride_id", http.StatusBadRequest)
		return
	}

	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	rideIDInt, err := strconv.Atoi(rideID)
	if err != nil {
		http.Error(w, "invalid ride_id", http.StatusBadRequest)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	ride, err := db.GetRideByID(d, rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting ride: %s", error), status)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if userIDInt != userAuth.UserID && ride.OwnerID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	if _, err := db.GetPassengerByRideIDAndUserID(d, rideIDInt, userIDInt); err != nil {
		log.WithError(err).Error("getting ride passenger")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting ride passenger: %s", error), status)
		return
	}

	if err := db.DeletePassenger(d, rideIDInt, userIDInt); err != nil {
		log.WithError(err).Error("deleting ride passenger")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   Feedback Handlers
// =====================

func getFeedbacks(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	feedbacks, err := db.GetFeedbacks(d)
	if err != nil {
		log.WithError(err).Error("getting feedbacks")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedbacks)
}

func getFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["feedback_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	feedback, err := db.GetFeedbackByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting feedback")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedback)
}

func getRideFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]

	if rideID == "" {
		http.Error(w, "missing ride_id", http.StatusBadRequest)
		return
	}

	rideIDInt, err := strconv.Atoi(rideID)
	if err != nil {
		http.Error(w, "invalid ride_id", http.StatusBadRequest)
		return
	}

	feedbacks, err := db.GetFeedbacksByRideID(d, rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting ride feedbacks")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedbacks)
}

func getUserFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if userIDInt != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	feedbacks, err := db.GetFeedbacksByUserID(d, userIDInt)
	if err != nil {
		log.WithError(err).Error("getting user feedbacks")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedbacks)
}

func getUserRideFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	rideID := vars["ride_id"]

	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	if rideID == "" {
		http.Error(w, "missing ride_id", http.StatusBadRequest)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	rideIDInt, err := strconv.Atoi(rideID)
	if err != nil {
		http.Error(w, "invalid ride_id", http.StatusBadRequest)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if userIDInt != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	feedbacks, err := db.GetFeedbackByUserIDAndRideID(d, userIDInt, rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting user ride feedback")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedbacks)
}

func createFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	var feedback core.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ride, err := db.GetRideByID(d, feedback.RideID)
	if err != nil {
		log.WithError(err).Error("getting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting ride: %s", error), status)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if feedback.UserID != userAuth.UserID && ride.OwnerID == userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	_, err = db.GetUserByID(d, int64(feedback.UserID))
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting user: %s", error), status)
		return
	}

	ride, err = db.GetRideByID(d, feedback.RideID)
	if err != nil {
		log.WithError(err).Error("getting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting ride: %s", error), status)
		return
	}

	passengers, err := db.GetPassengersByRideID(d, feedback.RideID)
	if err != nil {
		log.WithError(err).Error("getting ride passengers")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting ride passengers: %s", error), status)
		return
	}

	if err := feedback.Validate(ride, passengers); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := db.CreateFeedback(d, feedback)
	if err != nil {
		log.WithError(err).Error("creating feedback")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	feedback.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, feedback)
}

func updateFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["feedback_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var feedback core.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	existingFeedback, err := db.GetFeedbackByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting feedback")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting feedback: %s", error), status)
		return
	}

	// disallow changing the user or ride id
	feedback.UserID = existingFeedback.UserID
	feedback.RideID = existingFeedback.RideID

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if feedback.UserID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	ride, err := db.GetRideByID(d, feedback.RideID)
	if err != nil {
		log.WithError(err).Error("getting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting ride: %s", error), status)
		return
	}

	passengers, err := db.GetPassengersByRideID(d, feedback.RideID)
	if err != nil {
		log.WithError(err).Error("getting ride passengers")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting ride passengers: %s", error), status)
		return
	}

	if err := feedback.Validate(ride, passengers); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := db.UpdateFeedback(d, idInt, feedback); err != nil {
		log.WithError(err).Error("updating feedback")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["feedback_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	feedback, err := db.GetFeedbackByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting feedback")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting feedback: %s", error), status)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if feedback.UserID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	if err := db.DeleteFeedback(d, idInt); err != nil {
		log.WithError(err).Error("deleting feedback")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respond(w http.ResponseWriter, r *http.Request, data any) {
	acceptHeader := r.Header.Get("Accept")

	var response []byte
	var err error

	switch acceptHeader {
	case responseTypeXML:
		w.Header().Set("Content-Type", "application/xml")
		response, err = xml.Marshal(data)
	default:
		w.Header().Set("Content-Type", "application/json")
		response, err = json.Marshal(data)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
