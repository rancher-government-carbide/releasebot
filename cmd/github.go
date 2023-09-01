package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
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
		t.Time = time.Time{}
		log.Printf("Failed to decode into type time.Time")
		// temporary removal of error? future changes to project should render published_at timestamp irrelevant
		// return err
	}
	t.Time = parsedTime
	return nil
}

// sorts releases by publish date (newest to oldest)
func sortByPublishDate(releases []Release) []Release {
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].PublishedAt.After(releases[j].PublishedAt.Time)
	})
	return releases
}

// fetches only the single latest release from the repo (different endpoint than all releases)
//
//lint:ignore U1000 Ignore this unused function
func getLatestRelease(owner string, repo string) (Release, error) {

	var latestRelease Release

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Print("No provided github token - requests to the github api will be unathenticated (60 requests/hr rate limit)\n")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/%s/releases/latest", github_api_url, owner, repo), nil)
	if err != nil {
		log.Print(err)
		return latestRelease, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		log.Printf("Failed to get releases for: %s/%s", owner, repo)
		return latestRelease, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&latestRelease)
	if err != nil {
		log.Print(err)
		log.Printf("Failed to decode releases for: %s/%s", owner, repo)
		return latestRelease, err
	}
	return latestRelease, nil
}

// fetches all releases/prereleases for a repo (default gh api pagination is 30 results)
func getAllReleases(owner string, repo string) ([]Release, error) {
	var releases []Release
	token := os.Getenv("GITHUB_TOKEN")
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
	return releases, nil
}

// filters all prereleases out of the array (leaves only releases)
func filterPrereleases(releases []Release) []Release {
	var onlyRegularReleases []Release
	for _, release := range releases {
		if !release.Prerelease {
			onlyRegularReleases = append(onlyRegularReleases, release)
		}
	}
	return onlyRegularReleases
}

// filters all releases out of the array (leaves only prereleases)
func filterReleases(releases []Release) []Release {
	var onlyPrereleases []Release
	for _, release := range releases {
		if release.Prerelease {
			onlyPrereleases = append(onlyPrereleases, release)
		}
	}
	return onlyPrereleases
}

// fetches all releases from repo (sorted by publish date)
//
// count specifies the maximum number of releases to return, if its less than 0 there is no max
func getLatestReleases(owner string, repo string, prerelease bool, count int) ([]Release, error) {
	var latestReleases []Release
	latestReleases, err := getAllReleases(owner, repo)
	if err != nil {
		return latestReleases, err
	}
	if prerelease {
		latestReleases = filterReleases(latestReleases)
	} else {
		latestReleases = filterPrereleases(latestReleases)
	}
	// github api is generally already sorted by date already but they don't officially guarantee such
	latestReleases = sortByPublishDate(latestReleases)
	if count < 0 {
		return latestReleases, nil
	}
	if len(latestReleases) > (count - 1) {
		latestReleases = latestReleases[:count]
	}
	return latestReleases, nil
}
