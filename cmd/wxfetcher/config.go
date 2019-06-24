package main

import (
	"encoding/json"
	"io"
	"os"
)

type appConfig struct {
	DBConfig struct {
		Driver string `json:"driver"`
		Source string `json:"source"`
	} `json:"db"`
	LoggingConfig struct {
		Level string `json:"level"`
	} `json:"logging"`
}

func readConfig(path string) (cfg *appConfig, err error) {
	var f io.ReadCloser
	if f, err = os.Open(path); err != nil {
		return
	}
	defer f.Close()
	if err = json.NewDecoder(f).Decode(&cfg); err != nil {
		return
	}
	return
}
