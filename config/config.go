package config

import (
	"encoding/json"
	"os"
)

var (
	Configuration = loadConfiguration()
)

type SrvConfiguration struct {
	Blocklists []string `json:"Blocklists"`
	Reciever   string   `json:"Reciever"`
	Pusher     string   `json:"Pusher"`
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
