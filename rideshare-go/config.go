package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type config struct {
	MySQL  DBConfig `json:"db"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

func (c *DBConfig) dBConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
}

func loadConfig(path string) (*config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config config
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
