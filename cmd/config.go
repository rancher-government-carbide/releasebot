package main

import (
	"encoding/json"
	"log"
	"os"
)

// import (
//
// )

type Repository struct {
	Owner       string `json:"owner"`
	Repo        string `json:"repo"`
	Prereleases bool   `json:"prereleases"`
	Tekton      bool   `json:"tekton"`
	Slack       bool   `json:"slack"`
}

func loadConfig(config *[]Repository) error {

	configpath := os.Getenv("releasebot_config")
	if configpath == "" {
		log.Printf("Defaulting to ./config.json since no config path was specified.\n")
		configpath = "config.json"
	}

	configFile, err := os.Open(configpath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(config)
	if err != nil {
		return err
	}

	return nil
}
