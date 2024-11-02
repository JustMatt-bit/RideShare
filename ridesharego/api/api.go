package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"main/core"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const responseTypeXML = "application/xml"

func CreateRouter(db *sql.DB, authSecret string) *mux.Router {
	r := mux.NewRouter()
	r.Use(loggerMiddleware)

	dbMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), core.CtxDB, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	r.Use(dbMiddleware)

	// Google OAuth2 endpoint
	api := r.PathPrefix("/api/v1").Subrouter()

	withAdmin := func(handler func(http.ResponseWriter, *http.Request, *logrus.Entry, *sql.DB)) http.HandlerFunc {
		return authMiddleware(authSecret, core.RoleAdmin, withMiddleware(handler))
	}

	withUser := func(handler func(http.ResponseWriter, *http.Request, *logrus.Entry, *sql.DB)) http.HandlerFunc {
		return authMiddleware(authSecret, core.RoleUser, withMiddleware(handler))
	}

	withGuest := func(handler func(http.ResponseWriter, *http.Request, *logrus.Entry, *sql.DB)) http.HandlerFunc {
		return withMiddleware(handler)
	}

	// Ride endpoints
	api.HandleFunc("/rides", withGuest(getRides)).Methods("GET")
	api.HandleFunc("/ride/{ride_id}", withGuest(getRide)).Methods("GET")
	api.HandleFunc("/ride", withUser(createRide)).Methods("POST")
	api.HandleFunc("/ride/{ride_id}", withUser(updateRide)).Methods("PUT")
	api.HandleFunc("/ride/{ride_id}", withUser(deleteRide)).Methods("DELETE")
	api.HandleFunc("/user/{user_id}/rides", withUser(getUserRides)).Methods("GET")

	// User endpoints
	api.HandleFunc("/users", withAdmin(getUsers)).Methods("GET")
	api.HandleFunc("/user/{user_id}", withUser(getUser)).Methods("GET")
	api.HandleFunc("/user", withAdmin(createUser)).Methods("POST")
	api.HandleFunc("/user/{user_id}", withUser(updateUser)).Methods("PUT")
	api.HandleFunc("/user/{user_id}", withUser(deleteUser)).Methods("DELETE")

	// Car endpoints
	api.HandleFunc("/cars", withAdmin(getCars)).Methods("GET")
	api.HandleFunc("/car/{car_id}", withUser(getCar)).Methods("GET")
	api.HandleFunc("/car", withUser(createCar)).Methods("POST")
	api.HandleFunc("/car/{car_id}", withUser(updateCar)).Methods("PUT")
	api.HandleFunc("/car/{car_id}", withUser(deleteCar)).Methods("DELETE")
	api.HandleFunc("/user/{user_id}/cars", withUser(getUserCars)).Methods("GET")

	// Car model endpoints
	api.HandleFunc("/car_models", withGuest(getCarModels)).Methods("GET")
	api.HandleFunc("/car_model/{model_id}", withGuest(getCarModel)).Methods("GET")
	api.HandleFunc("/car_model", withAdmin(createCarModel)).Methods("POST")
	api.HandleFunc("/car_model/{model_id}", withAdmin(updateCarModel)).Methods("PUT")
	api.HandleFunc("/car_model/{model_id}", withAdmin(deleteCarModel)).Methods("DELETE")

	// Car make endpoints
	api.HandleFunc("/car_makes", withGuest(getCarMakes)).Methods("GET")
	api.HandleFunc("/car_make/{make_id}", withGuest(getCarMake)).Methods("GET")
	api.HandleFunc("/car_make", withAdmin(createCarMake)).Methods("POST")
	api.HandleFunc("/car_make/{make_id}", withAdmin(updateCarMake)).Methods("PUT")
	api.HandleFunc("/car_make/{make_id}", withAdmin(deleteCarMake)).Methods("DELETE")

	// Car category endpoints
	api.HandleFunc("/car_categories", withGuest(getCarCategories)).Methods("GET")
	api.HandleFunc("/car_category/{category_id}", withGuest(getCarCategory)).Methods("GET")
	api.HandleFunc("/car_category", withAdmin(createCarCategory)).Methods("POST")
	api.HandleFunc("/car_category/{category_id}", withAdmin(updateCarCategory)).Methods("PUT")
	api.HandleFunc("/car_category/{category_id}", withAdmin(deleteCarCategory)).Methods("DELETE")

	// Ride passenger endpoints
	api.HandleFunc("/ride/{ride_id}/passengers", withUser(getRidePassengers)).Methods("GET")
	api.HandleFunc("/ride/{ride_id}/passenger/{user_id}", withUser(createRidePassenger)).Methods("POST")
	api.HandleFunc("/ride/{ride_id}/passenger/{user_id}", withUser(deleteRidePassenger)).Methods("DELETE")

	// Feedback endpoints
	api.HandleFunc("/feedback", withAdmin(getFeedbacks)).Methods("GET")
	api.HandleFunc("/feedback/{feedback_id}", withAdmin(getFeedback)).Methods("GET")
	api.HandleFunc("/feedback", withUser(createFeedback)).Methods("POST")
	api.HandleFunc("/feedback/{feedback_id}", withUser(updateFeedback)).Methods("PUT")
	api.HandleFunc("/feedback/{feedback_id}", withUser(deleteFeedback)).Methods("DELETE")
	api.HandleFunc("/ride/{ride_id}/feedback", withGuest(getRideFeedback)).Methods("GET")
	api.HandleFunc("/user/{user_id}/feedback", withUser(getUserFeedback)).Methods("GET")
	api.HandleFunc("/user/{user_id}/ride/{ride_id}/feedback", withUser(getUserRideFeedback)).Methods("GET")

	return r
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
