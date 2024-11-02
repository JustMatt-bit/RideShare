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

func getCarCategories(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	categories, err := db.GetCarCategories(d)
	if err != nil {
		log.WithError(err).Error("getting car categories")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, categories)
}

func getCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["category_id"]

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

	category, err := db.GetCarCategoryByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting car category")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, category)
}

func createCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	var category core.CarCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		log.WithError(err).Error("decoding request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := category.Validate(); err != nil {
		log.WithError(err).Error("validating request")
		http.Error(w, fmt.Sprintf("validating request: %s", err.Error()), http.StatusBadRequest)
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

func updateCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["category_id"]

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

	var category core.CarCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		log.WithError(err).Error("decoding request")
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := category.Validate(); err != nil {
		log.WithError(err).Error("validating request")
		http.Error(w, fmt.Sprintf("validating request: %s", err.Error()), http.StatusBadRequest)
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

func deleteCarCategory(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["category_id"]

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
