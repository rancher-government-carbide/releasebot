package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

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

func filter_prerelease(releases []Release) []Release {
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

func get_latest_release(owner string, repo string) Release {

	token := os.Getenv("github_token")
	if token == "" {
		log.Print("No provided github token - requests to the github api will be unathenticated (60 requests/hr rate limit)\n")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo), nil)
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

	var latestrelease Release
	err = json.NewDecoder(resp.Body).Decode(&latestrelease)
	if err != nil {
		log.Print(err)
		log.Printf("Failed to decode releases for: %s/%s", owner, repo)
	}

	return latestrelease
}

func get_releases(owner string, repo string) []Release {

	token := os.Getenv("github_token")
	if token == "" {
		log.Print("No provided github token - requests to the github api will be unathenticated (60 requests/hr rate limit)\n")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo), nil)
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

	releases = filter_prerelease(releases)

	return releases
}

func releasecheck(owner string, repo string) {

	latestrelease := get_latest_release(owner, repo)
	loaded_release := &latestrelease
	fmt.Printf("Base release for %s/%s is %s\n", owner, repo, loaded_release.Name)
	for {
		latestrelease := get_latest_release(owner, repo)
		if latestrelease.Name != loaded_release.Name {
			*loaded_release = latestrelease
			fmt.Printf("Found a new release for %s/%s\n", owner, repo)
			slacknotif(*loaded_release, owner, repo)
		} else {
			fmt.Printf("No new releases for %s/%s\n", owner, repo)
		}
		time.Sleep(5 * time.Minute)
	}

}
