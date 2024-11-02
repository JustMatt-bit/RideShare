package core

import (
	"encoding/json"
	"errors"
	"time"
)

type UserAuth struct {
	Role   string
	UserID int
}

type Ride struct {
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

func (r *Ride) Validate(car *Car) error {
	if r.OwnerID == 0 {
		return errors.New("missing owner_user_id")
	}

	if r.VehicleID == 0 {
		return errors.New("missing vehicle_id")
	}

	if car.UserID != r.OwnerID {
		return errors.New("vehicle does not belong to owner")
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

type User struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	Role      string          `json:"role"`
	Settings  json.RawMessage `json:"settings"`
	CreatedAt string          `json:"created_at,omitempty"`
}

func (u *User) Validate() error {
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

type Car struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	LicensePlate string `json:"license_plate"`
	Year         int    `json:"year"`
	ModelID      int    `json:"model_id"`
}

func (c *Car) Validate(cars []Car) error {
	if c.LicensePlate == "" {
		return errors.New("missing license_plate")
	}

	if c.UserID == 0 {
		return errors.New("missing user_id")
	}

	for _, car := range cars {
		if car.LicensePlate == c.LicensePlate {
			return errors.New("license_plate already exists")
		}
	}

	if c.Year == 0 {
		return errors.New("missing year")
	}

	if c.ModelID == 0 {
		return errors.New("missing model_id")
	}

	return nil
}

type CarModel struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	MakeID     int    `json:"make_id"`
	Name       string `json:"name"`
}

func (cm *CarModel) Validate() error {
	if cm.CategoryID == 0 {
		return errors.New("missing category_id")
	}

	if cm.MakeID == 0 {
		return errors.New("missing make_id")
	}

	if cm.Name == "" {
		return errors.New("missing name")
	}

	return nil
}

type CarMake struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (cm *CarMake) Validate() error {
	if cm.Name == "" {
		return errors.New("missing name")
	}

	return nil
}

type CarCategory struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	PassengerCount int    `json:"passenger_count"`
}

func (cc *CarCategory) Validate() error {
	if cc.Name == "" {
		return errors.New("missing name")
	}

	if cc.PassengerCount == 0 {
		return errors.New("missing passenger_count")
	}

	return nil
}

type Passenger struct {
	RideID      int    `json:"ride_id"`
	PassengerID int    `json:"passenger_id"`
	CreatedAt   string `json:"created_at"`
}

func (p *Passenger) Validate(ride *Ride, passengers []Passenger, passengerCount int) error {
	if p.RideID == 0 {
		return errors.New("missing ride_id")
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", ride.StartDate)
	if err != nil {
		return err
	}

	if startTime.Before(time.Now()) {
		return errors.New("ride has already started")
	}

	for _, passenger := range passengers {
		if passenger.PassengerID == p.PassengerID {
			return errors.New("user is already a passenger")
		}
	}

	if len(passengers) >= passengerCount {
		return errors.New("ride is full")
	}

	if p.PassengerID == 0 {
		return errors.New("missing user_id")
	}

	return nil
}

type Feedback struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	RideID    int    `json:"ride_id"`
	Score     int    `json:"score"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

func (f *Feedback) Validate(ride *Ride, passengers []Passenger, role string) error {
	if f.UserID == 0 {
		return errors.New("missing user_id")
	}

	if f.RideID == 0 {
		return errors.New("missing ride_id")
	}

	if ride.OwnerID == f.UserID {
		return errors.New("owner cannot leave feedback")
	}

	passengerFound := false
	for _, passenger := range passengers {
		if passenger.PassengerID == f.UserID {
			passengerFound = true
			break
		}
	}

	if !passengerFound && role != RoleAdmin {
		return errors.New("user is not a passenger")
	}

	if f.Score < 1 || f.Score > 5 {
		return errors.New("invalid score")
	}

	if f.Message == "" {
		return errors.New("missing message")
	}

	return nil
}

type UserAuthRecord struct {
	UserID  int64  `json:"user_id"`
	Service string `json:"service"`
	Token   string `json:"token"`
}
