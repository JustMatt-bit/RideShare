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

func getFeedbacks(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	feedbacks, err := db.GetFeedbacks(d)
	if err != nil {
		log.WithError(err).Error("getting feedbacks")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedbacks)
}

func getFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["feedback_id"]

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

	feedback, err := db.GetFeedbackByID(d, idInt)
	if err != nil {
		log.WithError(err).Error("getting feedback")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedback)
}

func getRideFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
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

	feedbacks, err := db.GetFeedbacksByRideID(d, rideIDInt)
	if err != nil {
		log.WithError(err).Error("getting ride feedbacks")
		error, status := db.SqlErrorToHTTP(err)
		http.Error(w, error, status)
		return
	}

	respond(w, r, feedbacks)
}

func getUserFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
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

func getUserRideFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
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
		log.WithError(err).Error("parsing user_id")
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	rideIDInt, err := strconv.Atoi(rideID)
	if err != nil {
		log.WithError(err).Error("parsing ride_id")
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

func createFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	var feedback core.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		log.WithError(err).Error("decoding request")
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

	if err := feedback.Validate(ride, passengers, userAuth.Role); err != nil {
		log.WithError(err).Error("validating feedback")
		http.Error(w, fmt.Sprintf("validating request: %s", err.Error()), http.StatusBadRequest)
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

func updateFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["feedback_id"]

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

	var feedback core.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		log.WithError(err).Error("decoding request")
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

	if err := feedback.Validate(ride, passengers, userAuth.Role); err != nil {
		log.WithError(err).Error("validating feedback")
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

func deleteFeedback(w http.ResponseWriter, r *http.Request, log *logrus.Entry, d *sql.DB) {
	vars := mux.Vars(r)
	id := vars["feedback_id"]

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
