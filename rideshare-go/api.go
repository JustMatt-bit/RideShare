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
	api.HandleFunc("/ride", withLog(getRides)).Methods("GET")
	api.HandleFunc("/ride/{id}", withLog(getRide)).Methods("GET")
	api.HandleFunc("/ride", withLog(createRide)).Methods("POST")
	api.HandleFunc("/ride/{id}", withLog(updateRide)).Methods("PUT")
	api.HandleFunc("/ride/{id}", withLog(deleteRide)).Methods("DELETE")

	// User endpoints
	api.HandleFunc("/user", withLog(getUsers)).Methods("GET")
	api.HandleFunc("/user/{id}", withLog(getUser)).Methods("GET")
	api.HandleFunc("/user", withLog(createUser)).Methods("POST")
	api.HandleFunc("/user/{id}", withLog(updateUser)).Methods("PUT")
	api.HandleFunc("/user/{id}", withLog(deleteUser)).Methods("DELETE")

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
	id := vars["id"]

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
	id := vars["id"]

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
		http.Error(w, fmt.Sprintf("validating request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := dbGetRideByID(idInt); err != nil {
		log.WithError(err).Error("getting ride")
		error, status := sqlErrorToHTTP(err)
		http.Error(w, fmt.Sprintf("getting ride: %s", error), status)
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
	id := vars["id"]

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
		http.Error(w, fmt.Sprintf("getting ride: %s", error), status)
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
	id := vars["id"]

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
	id := vars["id"]

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
	id := vars["id"]

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
