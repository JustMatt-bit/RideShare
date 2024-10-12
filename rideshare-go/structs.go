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
	CreatedAt    string `json:"created_at"`
}

func (r *ride) validate() error {
	if r.OwnerID == 0 {
		return errors.New("missing owner_user_id")
	}
	if r.VehicleID == 0 {
		return errors.New("missing vehicle_id")
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
	CreatedAt string          `json:"created_at"`
}
