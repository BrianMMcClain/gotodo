package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func parseConfig(configPath string) (*Config, error) {
	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
