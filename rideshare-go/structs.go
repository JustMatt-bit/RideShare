package main

import (
	"encoding/json"
	"errors"
)

type ride struct {
	ID           int    `json:"id"`
	OwnerID      int    `json:"owner_user_id"`
	VehicleID    int    `json:"vehicle_id"`
	StartDate    string `json:"start_date"`
	StartCity    string `json:"start_city"`
	StartAddress string `json:"start_address"`
	EndCity      string `json:"end_city"`
	EndAddress   string `json:"end_address"`
	CreatedAt    string `json:"created_at,omitempty"`
}

func (r *ride) validate() error {
	if r.OwnerID == 0 {
		return errors.New("missing owner_user_id")
	}
	_, err := dbGetUserByID(r.OwnerID)
	if err != nil {
		return errors.New("user not found")
	}

	if r.VehicleID == 0 {
		return errors.New("missing vehicle_id")
	}
	_, err = dbGetCarByID(r.VehicleID)
	if err != nil {
		return errors.New("car not found")
	}

	if r.StartDate == "" {
		return errors.New("missing start_date")
	}
	if r.StartCity == "" {
		return errors.New("missing start_city")
	}
	if r.StartAddress == "" {
		return errors.New("missing start_address")
	}
	if r.EndCity == "" {
		return errors.New("missing end_city")
	}
	if r.EndAddress == "" {
		return errors.New("missing end_address")
	}
	return nil
}

type user struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	Settings  json.RawMessage `json:"settings"`
	CreatedAt string          `json:"created_at,omitempty"`
}

func (u *user) validate() error {
	if u.Name == "" {
		return errors.New("missing name")
	}
	if u.Email == "" {
		return errors.New("missing email")
	}
	if u.Password == "" {
		return errors.New("missing password")
	}
	return nil
}

type car struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	LicensePlate string `json:"license_plate"`
	Year         int    `json:"year"`
	ModelID      int    `json:"model_id"`
}

func (c *car) validate() error {
	if c.LicensePlate == "" {
		return errors.New("missing license_plate")
	}

	if c.UserID == 0 {
		return errors.New("missing user_id")
	}
	_, err := dbGetUserByID(c.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if c.Year == 0 {
		return errors.New("missing year")
	}
	if c.ModelID == 0 {
		return errors.New("missing model_id")
	}
	_, err = dbGetCarModelByID(c.ModelID)
	if err != nil {
		return errors.New("model not found")
	}
	return nil
}

type carModel struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	MakeID     int    `json:"make_id"`
	Name       string `json:"name"`
}

func (cm *carModel) validate() error {
	if cm.CategoryID == 0 {
		return errors.New("missing category_id")
	}
	_, err := dbGetCarCategoryByID(cm.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	if cm.MakeID == 0 {
		return errors.New("missing make_id")
	}
	_, err = dbGetCarMakeByID(cm.MakeID)
	if err != nil {
		return errors.New("make not found")
	}

	if cm.Name == "" {
		return errors.New("missing name")
	}
	return nil
}

type carMake struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (cm *carMake) validate() error {
	if cm.Name == "" {
		return errors.New("missing name")
	}
	return nil
}

type carCategory struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	PassengerCount int    `json:"passenger_count"`
}

func (cc *carCategory) validate() error {
	if cc.Name == "" {
		return errors.New("missing name")
	}
	if cc.PassengerCount == 0 {
		return errors.New("missing passenger_count")
	}
	return nil
}

type passenger struct {
	RideID      int    `json:"ride_id"`
	PassengerID int    `json:"passenger_id"`
	CreatedAt   string `json:"created_at"`
}

func (p *passenger) validate() error {
	if p.RideID == 0 {
		return errors.New("missing ride_id")
	}
	_, err := dbGetRideByID(p.RideID)
	if err != nil {
		return errors.New("ride not found")
	}

	if p.PassengerID == 0 {
		return errors.New("missing user_id")
	}
	_, err = dbGetUserByID(p.PassengerID)
	if err != nil {
		return errors.New("user not found")
	}
	return nil
}

type feedback struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	RideID    int    `json:"ride_id"`
	Score     int    `json:"score"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

func (f *feedback) validate() error {
	if f.UserID == 0 {
		return errors.New("missing user_id")
	}
	_, err := dbGetUserByID(f.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if f.RideID == 0 {
		return errors.New("missing ride_id")
	}
	_, err = dbGetRideByID(f.RideID)
	if err != nil {
		return errors.New("ride not found")
	}

	if f.Score == 0 {
		return errors.New("missing score")
	}
	if f.Message == "" {
		return errors.New("missing message")
	}
	return nil
}
