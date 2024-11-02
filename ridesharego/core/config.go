package core

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

type GoogleAuthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackURL  string `json:"callback_url"`
}

type Config struct {
	MySQL  DBConfig `json:"db"`
	Server struct {
		Port       string `json:"port"`
		AuthSecret string `json:"auth_secret"`
	} `json:"server"`
	GoogleAuth GoogleAuthConfig `json:"google_auth"`
}

func (c *DBConfig) DBConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
