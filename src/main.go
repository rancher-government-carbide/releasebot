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
		go releasecheck(config[i].Owner, config[i].Repo)
	}

	for {
		time.Sleep(time.Minute)
	}

}
