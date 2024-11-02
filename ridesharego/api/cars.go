package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/core"
	"main/db"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func getCars(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	cars, err := db.GetCars(d)
	if err != nil {
		log.WithError(err).Error("getting cars")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, cars)
}

func getCar(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["car_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.WithError(err).Error("parsing id")
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

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if car.UserID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	respond(w, r, car)
}

func createCar(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	var car core.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		log.WithError(err).Error("decoding request")
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
		log.WithError(err).Error("validating request")
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

func updateCar(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["car_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.WithError(err).Error("parsing id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var car core.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		log.WithError(err).Error("decoding request")
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
		log.WithError(err).Error("validating request")
		http.Error(w, fmt.Sprintf("validating request: %s", err.Error()), http.StatusBadRequest)
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

func deleteCar(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["car_id"]

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.WithError(err).Error("parsing id")
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
