package main

import (
	"encoding/json"
	"os"
)

type appConfig struct {
	DBConfig struct {
		Driver string `json:"driver"`
		Source string `json:"source"`
	} `json:"db"`
}

func readConfig(path string) (*appConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg appConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
