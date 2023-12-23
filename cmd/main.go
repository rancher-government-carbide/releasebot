package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {

	var repos []RepositoryEntry
	if err := loadRepos(&repos); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalf("Error loading repos file")
	}

	var payloads []PayloadEntry
	if err := loadPayloads(&payloads); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalf("Error loading payloads file")
	}

	Monitor(repos, payloads)

}
