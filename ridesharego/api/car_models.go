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

func getCarModels(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	models, err := db.GetCarModels(d)
	if err != nil {
		log.WithError(err).Error("getting car models")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, models)
}

func getCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["model_id"]

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

	model, err := db.GetCarModelByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting car model")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, model)
}

func createCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	var model core.CarModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		log.WithError(err).Error("decoding request")
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
		log.WithError(err).Error("validating request")
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

func updateCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["model_id"]

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

	var model core.CarModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		log.WithError(err).Error("decoding request")
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
		log.WithError(err).Error("validating request")
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

func deleteCarModel(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["model_id"]

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
