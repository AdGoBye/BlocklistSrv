package main

import (
	"encoding/json"
	"os"
)

type SrvConfiguration struct {
	Blocklists []string `json:"Blocklists"`
}

func loadConfiguration() (config SrvConfiguration) {
	configBytes, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}
	return config
}
