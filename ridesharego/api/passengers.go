package api

import (
	"database/sql"
	"fmt"
	"main/core"
	"main/db"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func getRidePassengers(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]

	if rideID == "" {
		http.Error(w, "missing ride_id", http.StatusBadRequest)
		return
	}

	rideIDInt, err := strconv.Atoi(rideID)
	if err != nil {
		log.WithError(err).Error("parsing ride_id")
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

	if !passengerFound && userAuth.Role != core.RoleAdmin {
		http.Error(w, "user is not a passenger", http.StatusForbidden)
		return
	}

	respond(w, r, passengers)
}

func createRidePassenger(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
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
		log.WithError(err).Error("parsing ride_id")
		http.Error(w, "invalid ride_id", http.StatusBadRequest)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.WithError(err).Error("parsing user_id")
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
		log.WithError(err).Error("validating request")
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

func deleteRidePassenger(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
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
		log.WithError(err).Error("parsing ride_id")
		http.Error(w, "invalid ride_id", http.StatusBadRequest)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.WithError(err).Error("parsing user_id")
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
