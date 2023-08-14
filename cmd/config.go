package main

import (
	"encoding/json"
	"log"
	"os"
)

type RepositoryEntry struct {
	Owner       string     `json:"owner"`
	Repo        string     `json:"repo"`
	Prereleases bool       `json:"prereleases"`
	Payloads    PayloadMap `json:"payloads"`
	Slack       bool       `json:"slack"`
}

type PayloadMap map[string]bool

type PayloadEntry struct {
	Name    string          `json:"name"`
	Url     string          `json:"url"`
	Payload json.RawMessage `json:"payload"`
}

func (p *PayloadMap) UnmarshalJSON(data []byte) error {
	var payloadStrings []string
	if err := json.Unmarshal(data, &payloadStrings); err != nil {
		return err
	}

	payloadMap := make(map[string]bool)
	for _, payload := range payloadStrings {
		payloadMap[payload] = true
	}

	*p = payloadMap
	return nil
}

func loadRepos(config *[]RepositoryEntry) error {

	configpath := os.Getenv("RELEASEBOT_REPOS")
	if configpath == "" {
		log.Printf("Defaulting to ./repos.json since no config path was specified.\n")
		configpath = "repos.json"
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

func loadPayloads(config *[]PayloadEntry) error {

	configpath := os.Getenv("RELEASEBOT_PAYLOADS")
	if configpath == "" {
		log.Printf("Defaulting to ./payloads.json since no config path was specified.\n")
		configpath = "payloads.json"
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
