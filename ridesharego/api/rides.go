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

func getRides(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	rides, err := db.GetRides(d)
	if err != nil {
		log.WithError(err).Error("getting rides")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, rides)
}

func getRide(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["ride_id"]

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

	ride, err := db.GetRideByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting ride")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, ride)
}

func createRide(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	var ride core.Ride
	if err := json.NewDecoder(r.Body).Decode(&ride); err != nil {
		log.WithError(err).Error("decoding request")
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
		log.WithError(err).Error("validating request")
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

func updateRide(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["ride_id"]

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

	var ride core.Ride
	if err := json.NewDecoder(r.Body).Decode(&ride); err != nil {
		log.WithError(err).Error("decoding request")
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
		log.WithError(err).Error("validating request")
		http.Error(w, fmt.Sprintf("validating request: %s", err.Error()), http.StatusBadRequest)
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

func deleteRide(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["ride_id"]

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
