package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v55/github"
)

var persist, _ = strconv.ParseBool(os.Getenv("PERSIST"))

func Monitor(repos []RepositoryEntry, payloads []PayloadEntry) {
	if persist {
		if err := ensureDataFolder(DataFolderPath); err != nil {
			log.Fatalf("%v", err)
		}
	}
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

	repoName := fmt.Sprintf("%s/%s", repo.Owner, repo.Repo)
	releaseType := "release"
	if prereleases {
		releaseType = "prerelease"
	}

	interval, err := strconv.ParseInt(os.Getenv("interval"), 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"parsedInterval": interval,
		}).Info("Interval environment variable not set or invalid - defaulting to 5 minutes")
		interval = 5
	}
	intervalTime := time.Duration(interval) * time.Minute

	loadInitialReleases := loadReleasesFromGithub
	if persist {
		loadInitialReleases = loadReleasesFromFile
	}
LoadInitialReleases:
	loadedReleasesMap, err := loadInitialReleases(repo, prereleases)
	if err != nil {
		log.WithFields(log.Fields{
			"releaseType": releaseType,
			"repoName":    repoName,
			"error":       err,
		}).Error("Failed to retrieve initial releases, retrying...")
		time.Sleep(intervalTime)
		goto LoadInitialReleases
	}

	for {

		loadedReleasesStrings := stringifyLoadedReleases(loadedReleasesMap)
		log.WithFields(log.Fields{
			"releaseType": releaseType,
			"repoName":    repoName,
			"error":       err,
			"hashmap":     strings.Join(loadedReleasesStrings, ", "),
		}).Debug()

	LoadNewReleases:
		latestReleases, err := getLatestReleases(repo.Owner, repo.Repo, prereleases, -1)
		if err != nil {
			log.WithFields(log.Fields{
				"releaseType": releaseType,
				"repoName":    repoName,
				"error":       err,
			}).Error("Failed to get latest releases")
			time.Sleep(intervalTime)
			goto LoadNewReleases
		}

		newReleases := checkForNewReleases(latestReleases, loadedReleasesMap)
		if len(newReleases) == 0 {
			log.WithFields(log.Fields{
				"releaseType": releaseType,
				"repoName":    repoName,
			}).Info("No new releases")
		} else {
			for _, release := range newReleases {
				log.WithFields(log.Fields{
					"releaseType": releaseType,
					"repoName":    repoName,
					"release":     release.GetTagName(),
				}).Info("Found new release")
				errors := newReleaseActions(repo, release, payloads)
				if len(errors) != 0 {
					for _, err := range errors {
						log.WithFields(log.Fields{
							"releaseType": releaseType,
							"repoName":    repoName,
							"release":     release.GetTagName(),
							"actionError": err,
						}).Error("Action failed")
					}
				}
			}
		}

		time.Sleep(intervalTime)
	}
}

// collection of actions to take when a new release is found
func newReleaseActions(repo RepositoryEntry, release *github.RepositoryRelease, payloads []PayloadEntry) []error {
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
	if persist {
		err := writeReleaseToFile(release.GetTagName(), repo.Owner, repo.Repo)
		if err != nil {
			errors = append(errors, fmt.Errorf("error writing release to file: %v", err))
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Checks if any releases in the array are new. If there are some returns an array of the new ones
// along with the newest timestamp among them. The timestamp is unchanged from the input if there are no new releases.
func checkForNewReleases(latestReleases []*github.RepositoryRelease, loadedReleasesMap map[string]bool) []*github.RepositoryRelease {
	var newReleases []*github.RepositoryRelease
	for _, release := range latestReleases {
		if !loadedReleasesMap[release.GetTagName()] {
			newReleases = append(newReleases, release)
			loadedReleasesMap[release.GetTagName()] = true
		}
	}
	return newReleases
}

func loadReleasesFromFile(repo RepositoryEntry, prereleases bool) (map[string]bool, error) {
	releaseFile := fmt.Sprintf(ReleaseFileFormat, DataFolderPath, repo.Owner, repo.Repo)
	var releaseMap map[string]bool
	_, err := os.Stat(releaseFile)
	if err == nil {
		releaseMap, err = readMapFromFile(releaseFile)
		if err != nil {
			return nil, err
		}
	} else if os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"owner": repo.Owner,
			"repo":  repo.Repo,
		}).Info("Release history file doesn't exist, initializing such now...")
		releaseMap, err = loadReleasesFromGithub(repo, prereleases)
		if err != nil {
			return nil, err
		}
		err = writeMapToFile(releaseMap, releaseFile)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return releaseMap, nil
}

// loads initial batch of releases into a hashmap and returns such along with the latest release timestamp
func loadReleasesFromGithub(repo RepositoryEntry, prereleases bool) (map[string]bool, error) {
	loadedReleasesMap := make(map[string]bool)
	baseReleases, err := getLatestReleases(repo.Owner, repo.Repo, prereleases, -1)
	if err != nil {
		return loadedReleasesMap, err
	}
	for _, release := range baseReleases {
		loadedReleasesMap[release.GetTagName()] = true
	}
	return loadedReleasesMap, nil
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

func writeReleaseToFile(releaseTag string, repoOwner string, repoName string) error {
	releaseHistoryFile := fmt.Sprintf(ReleaseFileFormat, DataFolderPath, repoOwner, repoName)
	err := appendStringToFile(releaseTag, releaseHistoryFile)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"releaseTag":         releaseTag,
		"releaseHistoryFile": releaseHistoryFile,
	}).Info("Appended release tag to release history file")
	return nil
}
