package main

import (
	"log"
	"time"
)

func main() {

	var config []Repository
	if err := loadConfig(&config); err != nil {
		log.Fatal("Error loading config file...\nExiting...\n")
	}

	for i := 0; i < len(config); i++ {
		go monitor_repo(config[i].Owner, config[i].Repo, config[i].Prereleases, config[i].Tekton, config[i].Slack)
	}

	for {
		time.Sleep(time.Minute)
	}

}
