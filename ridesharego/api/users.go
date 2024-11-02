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

func getUsers(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	users, err := db.GetUsers(d)
	if err != nil {
		log.WithError(err).Error("getting users")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, users)
}

func getUser(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["user_id"]

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

	user, err := db.GetUserByID(d, int64(idInt))
	if err != nil {
		log.WithError(err).Error("getting user")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	userAuth := r.Context().Value(core.CtxAuth).(*core.UserAuth)
	if user.ID != userAuth.UserID && userAuth.Role != core.RoleAdmin {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	respond(w, r, user)
}

func createUser(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	var user core.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		log.WithError(err).Error("validating user")
		http.Error(w, fmt.Sprintf("validating user: %s", err.Error()), http.StatusBadRequest)
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

func updateUser(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["user_id"]

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

	var userUpdate core.User
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		log.WithError(err).Error("decoding request")
		http.Error(w, "decoding request", http.StatusBadRequest)
		return
	}

	if err := userUpdate.Validate(); err != nil {
		log.WithError(err).Error("validating user")
		http.Error(w, fmt.Sprintf("validating user: %s", err.Error()), http.StatusBadRequest)
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

func deleteUser(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["user_id"]

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

func getUserRides(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["user_id"]

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

	rides, err := db.GetRidesByUserID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting user rides")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, rides)
}

func getUserCars(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["user_id"]

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

	cars, err := db.GetCarsByUserID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting user rides")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	if len(cars) == 0 {
		http.Error(w, "no cars found", http.StatusNotFound)
		return
	}

	respond(w, r, cars)
}
