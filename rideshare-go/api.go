package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ctxKey string

const (
	ctxLog ctxKey = "logger"

	responseTypeXML = "application/xml"
)

func createRouter() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.Use(loggerMiddleware)

	// Ride endpoints
	api.HandleFunc("/rides", withLog(getRides)).Methods("GET")
	api.HandleFunc("/ride/{ride_id}", withLog(getRide)).Methods("GET")
	api.HandleFunc("/ride", withLog(createRide)).Methods("POST")
	api.HandleFunc("/ride/{ride_id}", withLog(updateRide)).Methods("PUT")
	api.HandleFunc("/ride/{ride_id}", withLog(deleteRide)).Methods("DELETE")
	api.HandleFunc("/user/{user_id}/rides", withLog(getUserRides)).Methods("GET")

	// User endpoints
	api.HandleFunc("/users", withLog(getUsers)).Methods("GET")
	api.HandleFunc("/user/{user_id}", withLog(getUser)).Methods("GET")
	api.HandleFunc("/user", withLog(createUser)).Methods("POST")
	api.HandleFunc("/user/{user_id}", withLog(updateUser)).Methods("PUT")
	api.HandleFunc("/user/{user_id}", withLog(deleteUser)).Methods("DELETE")

	// Car endpoints
	api.HandleFunc("/cars", withLog(getCars)).Methods("GET")
	api.HandleFunc("/car/{car_id}", withLog(getCar)).Methods("GET")
	api.HandleFunc("/car", withLog(createCar)).Methods("POST")
	api.HandleFunc("/car/{car_id}", withLog(updateCar)).Methods("PUT")
	api.HandleFunc("/car/{car_id}", withLog(deleteCar)).Methods("DELETE")

	// Car model endpoints
	api.HandleFunc("/car_models", withLog(getCarModels)).Methods("GET")
	api.HandleFunc("/car_model/{model_id}", withLog(getCarModel)).Methods("GET")
	api.HandleFunc("/car_model", withLog(createCarModel)).Methods("POST")
	api.HandleFunc("/car_model/{model_id}", withLog(updateCarModel)).Methods("PUT")
	api.HandleFunc("/car_model/{model_id}", withLog(deleteCarModel)).Methods("DELETE")

	// Car make endpoints
	api.HandleFunc("/car_makes", withLog(getCarMakes)).Methods("GET")
	api.HandleFunc("/car_make/{make_id}", withLog(getCarMake)).Methods("GET")
	api.HandleFunc("/car_make", withLog(createCarMake)).Methods("POST")
	api.HandleFunc("/car_make/{make_id}", withLog(updateCarMake)).Methods("PUT")
	api.HandleFunc("/car_make/{make_id}", withLog(deleteCarMake)).Methods("DELETE")

	// Car category endpoints
	api.HandleFunc("/car_categories", withLog(getCarCategories)).Methods("GET")
	api.HandleFunc("/car_category/{category_id}", withLog(getCarCategory)).Methods("GET")
	api.HandleFunc("/car_category", withLog(createCarCategory)).Methods("POST")
	api.HandleFunc("/car_category/{category_id}", withLog(updateCarCategory)).Methods("PUT")
	api.HandleFunc("/car_category/{category_id}", withLog(deleteCarCategory)).Methods("DELETE")

	// Ride passenger endpoints
	api.HandleFunc("/ride/{ride_id}/passengers", withLog(getRidePassengers)).Methods("GET")
	api.HandleFunc("/ride/{ride_id}/passenger/{user_id}", withLog(createRidePassenger)).Methods("POST")
	api.HandleFunc("/ride/{ride_id}/passenger/{user_id}", withLog(deleteRidePassenger)).Methods("DELETE")

	// Feedback endpoints
	api.HandleFunc("/feedback", withLog(getFeedbacks)).Methods("GET")
	api.HandleFunc("/feedback/{feedback_id}", withLog(getFeedback)).Methods("GET")
	api.HandleFunc("/feedback", withLog(createFeedback)).Methods("POST")
	api.HandleFunc("/feedback/{feedback_id}", withLog(updateFeedback)).Methods("PUT")
	api.HandleFunc("/feedback/{feedback_id}", withLog(deleteFeedback)).Methods("DELETE")
	api.HandleFunc("/ride/{ride_id}/feedback", withLog(getRideFeedback)).Methods("GET")
	api.HandleFunc("/user/{user_id}/feedback", withLog(getUserFeedback)).Methods("GET")
	api.HandleFunc("/user/{user_id}/ride/{ride_id}/feedback", withLog(getUserRideFeedback)).Methods("GET")

	return r
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := setupLogging()
		log = log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Logger

		ctx := context.WithValue(r.Context(), ctxLog, log)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withLog(handler func(http.ResponseWriter, *http.Request, *logrus.Logger)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value(ctxLog).(*logrus.Logger)
		handler(w, r, log)
	}
}

// =====================
//   Ride Handlers
// =====================

func getRides(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	rides, err := dbGetRides()
	if err != nil {
		log.WithError(err).Error("getting rides")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, rides)
}

func getRide(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	ride, err := dbGetRideByID(idInt)
	if err != nil {
		log.WithError(err).Error("getting ride")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, ride)
}

func createRide(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	var ride ride
	if err := json.NewDecoder(r.Body).Decode(&ride); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ride.validate(); err != nil {
		http.Error(w, fmt.Sprintf("validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := dbCreateRide(ride)
	if err != nil {
		log.WithError(err).Error("creating ride")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	ride.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, ride)
}

func updateRide(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	var ride ride
	if err := json.NewDecoder(r.Body).Decode(&ride); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := ride.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := dbGetRideByID(idInt); err != nil {
		log.WithError(err).Error("getting ride")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting ride: %s", error), status)
		return
	}

	if err := dbUpdateRide(idInt, ride); err != nil {
		log.WithError(err).Error("updating ride")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteRide(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	if _, err := dbGetRideByID(idInt); err != nil {
		log.WithError(err).Error("getting ride")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("while getting ride: %s", error), status)
		return
	}

	if err := dbDeleteRide(idInt); err != nil {
		log.WithError(err).Error("deleting ride")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   User Handlers
// =====================

func getUsers(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	users, err := dbGetUsers()
	if err != nil {
		log.WithError(err).Error("getting users")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, users)
}

func getUser(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	user, err := dbGetUserByID(idInt)
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, user)
}

func createUser(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	var user user
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := dbCreateUser(user)
	if err != nil {
		log.WithError(err).Error("creating user")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	user.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, user)
}

func updateUser(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	var user user
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := user.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := dbGetUserByID(idInt); err != nil {
		log.WithError(err).Error("getting user")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting user: %s", error), status)
		return
	}

	if err := dbUpdateUser(idInt, user); err != nil {
		log.WithError(err).Error("updating user")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteUser(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	if _, err := dbGetUserByID(idInt); err != nil {
		log.WithError(err).Error("getting user")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting user: %s", error), status)
		return
	}

	if err := dbDeleteUser(idInt); err != nil {
		log.WithError(err).Error("deleting user")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getUserRides(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	rides, err := dbGetRidesByUserID(idInt)
	if err != nil {
		log.WithError(err).Error("getting user rides")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, rides)
}

// =====================
//   Car Handlers
// =====================

func getCars(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	cars, err := dbGetCars()
	if err != nil {
		log.WithError(err).Error("getting cars")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, cars)
}

func getCar(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	car, err := dbGetCarByID(idInt)
	if err != nil {
		log.WithError(err).Error("getting car")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, car)
}

func createCar(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	var car car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := car.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := dbCreateCar(car)
	if err != nil {
		log.WithError(err).Error("creating car")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	car.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, car)
}

func updateCar(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	var car car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := car.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := dbGetCarByID(idInt); err != nil {
		log.WithError(err).Error("getting car")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car: %s", error), status)
		return
	}

	if err := dbUpdateCar(idInt, car); err != nil {
		log.WithError(err).Error("updating car")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteCar(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	if _, err := dbGetCarByID(idInt); err != nil {
		log.WithError(err).Error("getting car")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car: %s", error), status)
		return
	}

	if err := dbDeleteCar(idInt); err != nil {
		log.WithError(err).Error("deleting car")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   Car Model Handlers
// =====================

func getCarModels(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	models, err := dbGetCarModels()
	if err != nil {
		log.WithError(err).Error("getting car models")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, models)
}

func getCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	model, err := dbGetCarModelByID(idInt)
	if err != nil {
		log.WithError(err).Error("getting car model")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, model)
}

func createCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	var model carModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := model.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := dbCreateCarModel(model)
	if err != nil {
		log.WithError(err).Error("creating car model")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	model.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, model)
}

func updateCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	var model carModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := model.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := dbGetCarModelByID(idInt); err != nil {
		log.WithError(err).Error("getting car model")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car model: %s", error), status)
		return
	}

	if err := dbUpdateCarModel(idInt, model); err != nil {
		log.WithError(err).Error("updating car model")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	if _, err := dbGetCarModelByID(idInt); err != nil {
		log.WithError(err).Error("getting car model")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car model: %s", error), status)
		return
	}

	if err := dbDeleteCarModel(idInt); err != nil {
		log.WithError(err).Error("deleting car model")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   Car Make Handlers
// =====================

func getCarMakes(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	makes, err := dbGetCarMakes()
	if err != nil {
		log.WithError(err).Error("getting car makes")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, makes)
}

func getCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	make, err := dbGetCarMakeByID(idInt)
	if err != nil {
		log.WithError(err).Error("getting car make")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, make)
}

func createCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	var make carMake
	if err := json.NewDecoder(r.Body).Decode(&make); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := make.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := dbCreateCarMake(make)
	if err != nil {
		log.WithError(err).Error("creating car make")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	make.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, make)
}

func updateCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	var make carMake
	if err := json.NewDecoder(r.Body).Decode(&make); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := make.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := dbGetCarMakeByID(idInt); err != nil {
		log.WithError(err).Error("getting car make")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car make: %s", error), status)
		return
	}

	if err := dbUpdateCarMake(idInt, make); err != nil {
		log.WithError(err).Error("updating car make")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	if _, err := dbGetCarMakeByID(idInt); err != nil {
		log.WithError(err).Error("getting car make")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car make: %s", error), status)
		return
	}

	if err := dbDeleteCarMake(idInt); err != nil {
		log.WithError(err).Error("deleting car make")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   Car Category Handlers
// =====================

func getCarCategories(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	categories, err := dbGetCarCategories()
	if err != nil {
		log.WithError(err).Error("getting car categories")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, categories)
}

func getCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	category, err := dbGetCarCategoryByID(idInt)
	if err != nil {
		log.WithError(err).Error("getting car category")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, category)
}

func createCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	var category carCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := category.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := dbCreateCarCategory(category)
	if err != nil {
		log.WithError(err).Error("creating car category")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	category.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, category)
}

func updateCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	var category carCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := category.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := dbGetCarCategoryByID(idInt); err != nil {
		log.WithError(err).Error("getting car category")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car category: %s", error), status)
		return
	}

	if err := dbUpdateCarCategory(idInt, category); err != nil {
		log.WithError(err).Error("updating car category")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	if _, err := dbGetCarCategoryByID(idInt); err != nil {
		log.WithError(err).Error("getting car category")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting car category: %s", error), status)
		return
	}

	if err := dbDeleteCarCategory(idInt); err != nil {
		log.WithError(err).Error("deleting car category")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   Ride Passenger Handlers
// =====================

func getRidePassengers(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	passengers, err := dbGetPassengersByRideID(rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting ride passengers")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, passengers)
}

func createRidePassenger(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	passenger := passenger{
		RideID:      rideIDInt,
		PassengerID: userIDInt,
	}

	if err := passenger.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = dbCreatePassenger(passenger)
	if err != nil {
		log.WithError(err).Error("creating ride passenger")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respond(w, r, passenger)
}

func deleteRidePassenger(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	if _, err := dbGetPassengerByRideIDAndUserID(rideIDInt, userIDInt); err != nil {
		log.WithError(err).Error("getting ride passenger")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting ride passenger: %s", error), status)
		return
	}

	if err := dbDeletePassenger(rideIDInt, userIDInt); err != nil {
		log.WithError(err).Error("deleting ride passenger")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =====================
//   Feedback Handlers
// =====================

func getFeedbacks(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	feedbacks, err := dbGetFeedbacks()
	if err != nil {
		log.WithError(err).Error("getting feedbacks")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedbacks)
}

func getFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	feedback, err := dbGetFeedbackByID(idInt)
	if err != nil {
		log.WithError(err).Error("getting feedback")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedback)
}

func getRideFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	feedbacks, err := dbGetFeedbacksByRideID(rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting ride feedbacks")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedbacks)
}

func getUserFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	feedbacks, err := dbGetFeedbacksByUserID(userIDInt)
	if err != nil {
		log.WithError(err).Error("getting user feedbacks")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedbacks)
}

func getUserRideFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	feedbacks, err := dbGetFeedbackByUserIDAndRideID(userIDInt, rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting user ride feedback")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedbacks)
}

func createFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
	var feedback feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := feedback.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := dbCreateFeedback(feedback)
	if err != nil {
		log.WithError(err).Error("creating feedback")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}
	feedback.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	respond(w, r, feedback)
}

func updateFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	var feedback feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	existingFeedback, err := dbGetFeedbackByID(idInt)
	if err != nil {
		log.WithError(err).Error("getting feedback")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting feedback: %s", error), status)
		return
	}

	// disallow changing the user or ride id
	feedback.UserID = existingFeedback.UserID
	feedback.RideID = existingFeedback.RideID

	if err := feedback.validate(); err != nil {
		http.Error(w, fmt.Sprintf("while validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := dbUpdateFeedback(idInt, feedback); err != nil {
		log.WithError(err).Error("updating feedback")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Logger) {
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

	if _, err := dbGetFeedbackByID(idInt); err != nil {
		log.WithError(err).Error("getting feedback")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting feedback: %s", error), status)
		return
	}

	if err := dbDeleteFeedback(idInt); err != nil {
		log.WithError(err).Error("deleting feedback")
		error, status := sqlErrorToHTTP(err)
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
