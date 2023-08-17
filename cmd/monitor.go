package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func initMonitor(repos []RepositoryEntry, payloads []PayloadEntry) {
	for i := 0; i < len(repos); i++ {
		if repos[i].Prereleases {
			go monitorRepo(repos[i], payloads, true)
		}
		go monitorRepo(repos[i], payloads, false)
	}
	for {
		time.Sleep(time.Minute)
	}
}

// periodically checks the github api for new releases
func monitorRepo(repo RepositoryEntry, payloads []PayloadEntry, prereleases bool) {

	releaseType := "release"
	if prereleases {
		releaseType = "prerelease"
	}

	interval, err := strconv.ParseInt(os.Getenv("interval"), 10, 64)
	if err != nil {
		log.Printf("Interval environment variable not set or invalid value - defaulting to 5 minutes")
		interval = 5
	}

	loadedReleasesMap, newestReleaseTimestamp, err := loadInitialReleases(repo, prereleases)
	if err != nil {
		log.Printf("Failed to fetch initial %ss for %s/%s: %v", releaseType, repo.Owner, repo.Repo, err)
	}

	for {

		// temporary debugging: print the releases in the map
		loadedReleasesStrings := stringifyLoadedReleases(loadedReleasesMap)
		log.Printf("%ss in the hashmap for %s/%s: %s\n", releaseType, repo.Owner, repo.Repo, strings.Join(loadedReleasesStrings, ", "))

		latestReleases, err := getLatestReleases(repo.Owner, repo.Repo, prereleases, 10)
		if err != nil {
			log.Printf("Failed to get latest %ss for %s/%s: %v", releaseType, repo.Owner, repo.Repo, err)
		}

		newReleases, updatedNewestReleaseTimestamp := checkForNewReleases(latestReleases, loadedReleasesMap, newestReleaseTimestamp)
		if len(newReleases) == 0 {
			log.Printf("No new %ss for %s/%s\n", releaseType, repo.Owner, repo.Repo)
		} else {
			log.Printf("Found new %ss for %s/%s: %v\n", releaseType, repo.Owner, repo.Repo, newReleases)
		}
		newestReleaseTimestamp = updatedNewestReleaseTimestamp

		for _, release := range newReleases {
			errors := newReleaseActions(repo, release, payloads)
			if len(errors) != 0 {
				for _, err := range errors {
					log.Printf("Action failed for new %s %s/%s (%s): %v", releaseType, repo.Owner, repo.Repo, release.TagName, err)
				}
			}
		}

		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

// collection of actions to take when a new release is found
func newReleaseActions(repo RepositoryEntry, release Release, payloads []PayloadEntry) []error {
	var errors []error
	if repo.Slack {
		err := slacknotif(release, repo.Owner, repo.Repo)
		if err != nil {
			errors = append(errors, fmt.Errorf("error sending Slack notification: %v", err))
		}
	}
	err := sendAllPayloads(release, repo, payloads)
	if err != nil {
		errors = append(errors, fmt.Errorf("error sending payload: %v", err))
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Checks if any releases in the array are new. If there are some returns an array of the new ones
// along with the newest timestamp among them. The timestamp is unchanged from the input if there are no new releases.
func checkForNewReleases(latestReleases []Release, loadedReleasesMap map[string]bool, newestReleaseTimestamp time.Time) ([]Release, time.Time) {
	var updatedNewestReleaseTimestamp time.Time
	var newReleases []Release
	for _, release := range latestReleases {
		if !loadedReleasesMap[release.TagName] && release.PublishedAt.After(newestReleaseTimestamp) {
			newReleases = append(newReleases, release)
			loadedReleasesMap[release.TagName] = true
			updatedNewestReleaseTimestamp = release.PublishedAt.Time
		}
	}
	return newReleases, updatedNewestReleaseTimestamp
}

// loads initial batch of releases into a hashmap and returns such along with the latest release timestamp
func loadInitialReleases(repo RepositoryEntry, prereleases bool) (map[string]bool, time.Time, error) {
	var newestReleaseTimestamp time.Time
	loadedReleasesMap := make(map[string]bool)
	baseReleases, err := getLatestReleases(repo.Owner, repo.Repo, prereleases, -1)
	if err != nil {
		return loadedReleasesMap, time.Time{}, err
	}
	for _, release := range baseReleases {
		if release.PublishedAt.After(newestReleaseTimestamp) {
			newestReleaseTimestamp = release.PublishedAt.Time
		}
		if !loadedReleasesMap[release.TagName] {
			loadedReleasesMap[release.TagName] = true
		}
	}
	return loadedReleasesMap, newestReleaseTimestamp, nil
}

// takes the current release hashmap and returns an array of all the release names in such
func stringifyLoadedReleases(loadedReleasesMap map[string]bool) []string {
	var loadedReleasesMessage = make([]string, 0, len(loadedReleasesMap))
	for release := range loadedReleasesMap {
		if loadedReleasesMap[release] {
			loadedReleasesMessage = append(loadedReleasesMessage, release)
		}
	}
	return loadedReleasesMessage
}
