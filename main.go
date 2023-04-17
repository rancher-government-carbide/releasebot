package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Release struct {
	Name        string `json:"name"`
	TagName     string `json:"tag_name"`
	Prerelease  bool   `json:"prerelease"`
	PublishedAt string `json:"published_at"`
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

	var latestrelease Release
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		log.Printf("Failed to get releases for: %s/%s", owner, repo)
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&latestrelease)
	if err != nil {
		log.Print(err)
		log.Printf("Failed to decode releases for: %s/%s", owner, repo)
	}

	return latestrelease
}

func get_releases(owner string, repo string) []Release {

	var releases []Release
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		log.Printf("Failed to get releases for: %s/%s", owner, repo)
	}

	defer resp.Body.Close()
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
	fmt.Printf("Base release for %s/%s is %s with a memory address of %p\n", owner, repo, loaded_release.Name, &latestrelease)
	for {
		latestrelease := get_latest_release(owner, repo)
		if latestrelease.Name != loaded_release.Name {
			*loaded_release = latestrelease
			fmt.Printf("Found a new release for %s/%s\n", owner, repo)
		} else {
			fmt.Printf("No new releases for %s/%s\n", owner, repo)
		}
		time.Sleep(10 * time.Second)
	}

}

func main() {

	go releasecheck("k3s-io", "k3s")
	go releasecheck("clanktron", "dummy")

	for {
		time.Sleep(time.Minute)
	}
	// go releasecheck("rancher", "rancher")
	// releasecheck("kubewarden", "kubewarden-controller")

}
