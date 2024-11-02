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

func getCarMakes(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	makes, err := db.GetCarMakes(d)
	if err != nil {
		log.WithError(err).Error("getting car makes")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, makes)
}

func getCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["make_id"]

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

	make, err := db.GetCarMakeByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting car make")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, make)
}

func createCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	var make core.CarMake
	if err := json.NewDecoder(r.Body).Decode(&make); err != nil {
		log.WithError(err).Error("decoding request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := make.Validate(); err != nil {
		log.WithError(err).Error("validating request")
		http.Error(w, fmt.Sprintf("validating request: %s", err.Error()), http.StatusBadRequest)
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

func updateCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["make_id"]

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

	var make core.CarMake
	if err := json.NewDecoder(r.Body).Decode(&make); err != nil {
		log.WithError(err).Error("decoding request")
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := make.Validate(); err != nil {
		log.WithError(err).Error("validating request")
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

func deleteCarMake(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["make_id"]

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
