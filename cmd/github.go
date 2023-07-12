package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

var github_api_url string = "https://api.github.com"

type Release struct {
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	Prerelease  bool   `json:"prerelease"`
	HtmlUrl     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
	Author      Author `json:"author"`
}

type Author struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
}

func sort_by_published_date(releases []Release) []Release {

	sort.Slice(releases, func(i, j int) bool {
		// Parse the ISO 8601 timestamps into time.Time objects
		t1, err := time.Parse(time.RFC3339, releases[i].PublishedAt)
		if err != nil {
			return false
		}

		t2, err := time.Parse(time.RFC3339, releases[j].PublishedAt)
		if err != nil {
			return false
		}

		// Compare the timestamps
		return t1.After(t2)
	})

	return releases
}

func get_releases(owner string, repo string) ([]Release, error) {

	var releases []Release
	token := os.Getenv("github_token")
	if token == "" {
		log.Print("No provided github token - requests to the github api will be unathenticated (60 requests/hr rate limit)\n")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/%s/releases", github_api_url, owner, repo), nil)
	if err != nil {
		return releases, fmt.Errorf("Failed to create http request for %s/%s - %w", owner, repo, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return releases, fmt.Errorf("Failed to send http request for %s/%s - %w", owner, repo, err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return releases, fmt.Errorf("Failed to decode releases for %s/%s - %w", owner, repo, err)
	}

	// log.Printf(" releases for: %s/%s - %v", owner, repo, releases)

	return releases, nil
}

// filters all prereleases out of the array (leaves only releases)
func filter_prereleases(releases []Release) []Release {
	var filteredreleases []Release
	for _, release := range releases {
		if !release.Prerelease {
			filteredreleases = append(filteredreleases, release)
		}
	}
	return filteredreleases
}

// filters all releases out of the array (leaves only prereleases)
func filter_releases(releases []Release) []Release {
	var filteredreleases []Release
	for _, release := range releases {
		if release.Prerelease {
			filteredreleases = append(filteredreleases, release)
		}
	}
	return filteredreleases
}

// fetches all releases from repo and returns the 5 latest releases - can fetch either releases or prereleases
func get_latest_releases(owner string, repo string, prerelease bool) ([]Release, error) {

	var latest_releases []Release
	latest_releases, err := get_releases(owner, repo)
	if err != nil {
		return latest_releases, err
	}
	if prerelease {
		latest_releases = filter_releases(latest_releases)
	} else {
		latest_releases = filter_prereleases(latest_releases)
	}
	// github api is generally already sorted by date already but they don't officially guarantee such
	latest_releases = sort_by_published_date(latest_releases)
	// only the latest 5 are relevant
	// if len(latest_releases) > 5 {
	// 	latest_releases = latest_releases[:5]
	// }
	return latest_releases, nil
}

func monitor_repo(owner string, repo string, prereleases bool, tekton bool, slack bool) {

	interval, err := strconv.ParseInt(os.Getenv("interval"), 10, 64)
	if err != nil {
		log.Printf("Interval environment variable not set or invalid value - defaulting to 5 minutes")
		interval = 5
	}

	var new_releases []Release
	var new_prereleases []Release

	loaded_releases_map := make(map[string]bool)
	loaded_prereleases_map := make(map[string]bool)

	base_releases, err := get_latest_releases(owner, repo, false)
	if err != nil {
		log.Printf("Failed to get latest releases for: %s/%s - %v", owner, repo, err)
	}
	for _, release := range base_releases {
		if !loaded_releases_map[release.TagName] {
			loaded_releases_map[release.TagName] = true
		}
	}
	base_prereleases, err := get_latest_releases(owner, repo, true)
	if err != nil {
		log.Printf("Failed to get latest releases for: %s/%s - %v", owner, repo, err)
	}
	for _, release := range base_prereleases {
		if !loaded_prereleases_map[release.TagName] {
			loaded_prereleases_map[release.TagName] = true
		}
	}

	for {
		new_releases, err = get_latest_releases(owner, repo, false)
		if err != nil {
			log.Printf("Failed to get latest releases for: %s/%s - %v", owner, repo, err)
		}
		if check_releases(&loaded_releases_map, new_releases, owner, repo, tekton, slack, false) {
			log.Printf("No new releases for %s/%s\n", owner, repo)
		}
		if prereleases {
			new_prereleases, err = get_latest_releases(owner, repo, true)
			if err != nil {
				log.Printf("Failed to get latest prereleases for: %s/%s - %v", owner, repo, err)
			}
			if check_releases(&loaded_prereleases_map, new_prereleases, owner, repo, tekton, slack, true) {
				log.Printf("No new prereleases for %s/%s\n", owner, repo)
			}
		}
		time.Sleep(time.Duration(interval) * time.Minute)
	}

}

// returns false if theres a new release
func check_releases(loaded_release_map *map[string]bool, new_releases []Release, owner string, repo string, tekton bool, slack bool, prerelease bool) bool {
	release_type := "release"
	if prerelease {
		release_type = "prerelease"
	}
	no_new_releases := true
	for _, release := range new_releases {
		if !(*loaded_release_map)[release.TagName] {
			log.Printf("Found a new %s for %s/%s (%s)", release_type, owner, repo, release.TagName)
			(*loaded_release_map)[release.TagName] = true
			if slack {
				err := slacknotif(release, owner, repo)
				if err != nil {
					log.Printf("Error sending slack notification for %s %s/%s %s", release_type, owner, repo, release.TagName)
				}
			}
			if tekton {
				err := triggertekton(release, owner, repo)
				if err != nil {
					log.Printf("Error sending tekton payload for %s %s/%s %s", release_type, owner, repo, release.TagName)
				}
			}
			no_new_releases = false
		}
	}
	return no_new_releases
}

// fetches only the latest release from the repo (different endpoint than all releases)
func get_latest_release(owner string, repo string) (Release, error) {

	var latestrelease Release

	token := os.Getenv("github_token")
	if token == "" {
		log.Print("No provided github token - requests to the github api will be unathenticated (60 requests/hr rate limit)\n")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/%s/releases/latest", github_api_url, owner, repo), nil)
	if err != nil {
		log.Print(err)
		return latestrelease, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		log.Printf("Failed to get releases for: %s/%s", owner, repo)
		return latestrelease, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&latestrelease)
	if err != nil {
		log.Print(err)
		log.Printf("Failed to decode releases for: %s/%s", owner, repo)
		return latestrelease, err
	}

	return latestrelease, nil
}
