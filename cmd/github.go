package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var github_api_url string = "https://api.github.com"

type Release struct {
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	Prerelease  bool   `json:"prerelease"`
	HtmlUrl     string `json:"html_url"`
	PublishedAt Time   `json:"published_at"`
	Author      Author `json:"author"`
}

type Author struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	const layout = time.RFC3339
	parsedTime, err := time.Parse(layout, string(data[1:len(data)-1]))
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

// sorts releases by publish date (newest to oldest)
func sort_by_published_date(releases []Release) []Release {
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].PublishedAt.After(releases[j].PublishedAt.Time)
	})
	return releases
}

// fetches all releases/prereleases for a repo (default gh api max is 30 results)
func get_releases(owner string, repo string) ([]Release, error) {

	var releases []Release
	token := os.Getenv("github_token")
	if token == "" {
		log.Print("No provided github token - requests to the github api will be unathenticated (60 requests/hr rate limit)\n")
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/%s/releases", github_api_url, owner, repo), nil)
	if err != nil {
		return releases, fmt.Errorf("failed to create http request for %s/%s - %w", owner, repo, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return releases, fmt.Errorf("failed to send http request for %s/%s - %w", owner, repo, err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return releases, fmt.Errorf("failed to decode releases for %s/%s - %w", owner, repo, err)
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

// fetches all releases from repo (sorted by publish date)
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

	release_type := "release"
	if prereleases {
		release_type = "prerelease"
	}

	interval, err := strconv.ParseInt(os.Getenv("interval"), 10, 64)
	if err != nil {
		log.Printf("Interval environment variable not set or invalid value - defaulting to 5 minutes")
		interval = 5
	}

	var new_releases []Release
	var newest_release_published time.Time
	loaded_releases_map := make(map[string]bool)

	base_releases, err := get_latest_releases(owner, repo, prereleases)
	if err != nil {
		log.Printf("Failed to get latest %ss for: %s/%s - %v", release_type, owner, repo, err)
	}

	for _, release := range base_releases {
		if release.PublishedAt.After(newest_release_published) {
			newest_release_published = release.PublishedAt.Time
		}
		if !loaded_releases_map[release.TagName] {
			loaded_releases_map[release.TagName] = true
		}
	}

	for {
		// temporary debugging: print the releases in the map
		var loaded_releases_msg = make([]string, 0, len(loaded_releases_map))
		for release := range loaded_releases_map {
			loaded_releases_msg = append(loaded_releases_msg, release)
		}
		log.Printf("%ss in the hashmap for %s/%s: %s\n", release_type, owner, repo, strings.Join(loaded_releases_msg, ", "))

		new_releases, err = get_latest_releases(owner, repo, prereleases)
		if err != nil {
			log.Printf("Failed to get latest %ss for: %s/%s - %v", release_type, owner, repo, err)
		}
		if len(new_releases) > 4 {
			new_releases = new_releases[:5]
		}

		no_new_releases := true
		for _, release := range new_releases {
			if !loaded_releases_map[release.TagName] && release.PublishedAt.After(newest_release_published) {
				log.Printf("Found a new %s for %s/%s (%s)", release_type, owner, repo, release.TagName)
				loaded_releases_map[release.TagName] = true
				newest_release_published = release.PublishedAt.Time
				if slack {
					err := slacknotif(release, owner, repo)
					if err != nil {
						log.Printf("Error sending slack notification for %s %s/%s (%s) %v", release_type, owner, repo, release.TagName, err)
					}
				}
				if tekton {
					err := triggertekton(release, owner, repo)
					if err != nil {
						log.Printf("Error sending tekton payload for %s %s/%s (%s) %v", release_type, owner, repo, release.TagName, err)
					}
				}
				no_new_releases = false
			}
		}
		if no_new_releases {
			log.Printf("No new %ss for %s/%s\n", release_type, owner, repo)
		}

		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

// // fetches only the latest release from the repo (different endpoint than all releases)
// func get_latest_release(owner string, repo string) (Release, error) {
//
// 	var latestrelease Release
//
// 	token := os.Getenv("github_token")
// 	if token == "" {
// 		log.Print("No provided github token - requests to the github api will be unathenticated (60 requests/hr rate limit)\n")
// 	}
//
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/%s/releases/latest", github_api_url, owner, repo), nil)
// 	if err != nil {
// 		log.Print(err)
// 		return latestrelease, err
// 	}
// 	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Print(err)
// 		log.Printf("Failed to get releases for: %s/%s", owner, repo)
// 		return latestrelease, err
// 	}
// 	defer resp.Body.Close()
//
// 	err = json.NewDecoder(resp.Body).Decode(&latestrelease)
// 	if err != nil {
// 		log.Print(err)
// 		log.Printf("Failed to decode releases for: %s/%s", owner, repo)
// 		return latestrelease, err
// 	}
//
// 	return latestrelease, nil
// }
