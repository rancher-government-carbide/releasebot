package main

import (
	"log"
)

func main() {

	var repos []RepositoryEntry
	if err := loadRepos(&repos); err != nil {
		log.Fatalf("Error loading repos file: %v", err)
	}

	var payloads []PayloadEntry
	if err := loadPayloads(&payloads); err != nil {
		log.Fatalf("Error loading payloads file: %v", err)
	}

	Monitor(repos, payloads)

}
