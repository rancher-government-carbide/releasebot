package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const github_api_url string = "https://api.github.com"

type Release struct {
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	Prerelease  bool   `json:"prerelease"`
	HtmlUrl     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
	Author      struct {
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
		HtmlUrl   string `json:"html_url"`
	} `json:"author"`
}

func get_releases(owner string, repo string) []Release {

	token := os.Getenv("github_token")
	if token == "" {
		log.Print("No provided github token - requests to the github api will be unathenticated (60 requests/hr rate limit)\n")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/repos/%s/%s/releases", github_api_url, owner, repo), nil)
	if err != nil {
		log.Print(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		log.Printf("Failed to get releases for: %s/%s", owner, repo)
	}
	defer resp.Body.Close()

	var releases []Release
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		log.Print(err)
		log.Printf("Failed to decode releases for: %s/%s", owner, repo)
	}

	return releases
}

func filter_prereleases(releases []Release) []Release {
	var filteredreleases []Release
	for _, release := range releases {
		if !release.Prerelease {
			filteredreleases = append(filteredreleases, release)
		}
	}
	// only checking the first 10 results
	filteredreleases = filteredreleases[:10]
	return filteredreleases
}

func filter_releases(releases []Release) []Release {
	var filteredreleases []Release
	for _, release := range releases {
		if release.Prerelease {
			filteredreleases = append(filteredreleases, release)
		}
	}
	// only checking the first 10 results
	filteredreleases = filteredreleases[:10]
	return filteredreleases
}

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

func get_latest_prerelease(owner string, repo string) (Release, error) {

	var releases []Release
	releases = get_releases(owner, repo)
	releases = filter_releases(releases)

	return releases[0], nil
}

func monitor_repo(owner string, repo string, prereleases bool, tekton bool) {

	interval, err := strconv.ParseInt(os.Getenv("interval"), 10, 64)
	if err != nil {
		log.Printf("Interval environment variable not set or invalid value - defaulting to 5 minutes")
		interval = 5
	}

	var loaded_release Release
	var latest_release Release

	var loaded_prerelease Release
	var latest_prerelease Release

	loaded_release, err = get_latest_release(owner, repo)
	if err != nil {
		log.Printf("Using placeholder '' as current release for: %s/%s", owner, repo)
		loaded_release.Name = ""
	}
	fmt.Printf("Base release for %s/%s is %s\n", owner, repo, loaded_release.Name)

	loaded_prerelease, err = get_latest_prerelease(owner, repo)
	if err != nil {
		log.Printf("Using placeholder '' as current prerelease for: %s/%s", owner, repo)
		loaded_release.Name = ""
	}
	fmt.Printf("Base prerelease for %s/%s is %s\n", owner, repo, loaded_prerelease.Name)

	for {

		latest_release, err = get_latest_release(owner, repo)
		if err != nil {
			log.Printf("Failed to get latest release for: %s/%s", owner, repo)
		} else {
			if latest_release.Name != loaded_release.Name {
				loaded_release = latest_release
				fmt.Printf("Found a new release for %s/%s\n", owner, repo)
				slacknotif(loaded_release, owner, repo, releases_channel)
				if tekton {
					triggertekton(loaded_release, owner, repo)
				}
			} else {
				fmt.Printf("No new releases for %s/%s\n", owner, repo)
			}
		}

		if prereleases {

			latest_prerelease, err = get_latest_prerelease(owner, repo)
			if err != nil {
				log.Printf("Failed to get latest release for: %s/%s", owner, repo)
			} else {
				if latest_prerelease.Name != loaded_prerelease.Name {
					loaded_prerelease = latest_prerelease
					fmt.Printf("Found a new release for %s/%s\n", owner, repo)
					slacknotif(loaded_release, owner, repo, prerelease_channel)
					if tekton {
						triggertekton(loaded_release, owner, repo)
					}
				} else {
					fmt.Printf("No new releases for %s/%s\n", owner, repo)
				}
			}
		}

		time.Sleep(time.Duration(interval) * time.Minute)
	}

}
