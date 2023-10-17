package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"

	"github.com/google/go-github/v55/github"
)

// sorts releases by publish date (newest to oldest)
func sortByPublishDate(releases []*github.RepositoryRelease) []*github.RepositoryRelease {
	sort.Slice(releases, func(i, j int) bool {
		if releases[i].PublishedAt == nil && releases[j].PublishedAt == nil {
			return false // If both are nil, consider them equal
		} else if releases[i].PublishedAt == nil {
			return false // Nil is considered older than non-nil
		} else if releases[j].PublishedAt == nil {
			return true // Non-nil is considered newer than nil
		}
		return releases[i].PublishedAt.After(releases[j].PublishedAt.Time)
	})
	return releases
}

// fetches all releases/prereleases for a repo (default gh api pagination is 30 results)
func getAllReleases(owner string, repo string) ([]*github.RepositoryRelease, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Info("No provided github token - requests to the github api will be unathenticated (60 requests/hr rate limit)\n")
	}
	client := github.NewClient(nil).WithAuthToken(token)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 100}

	var allReleases []*github.RepositoryRelease
	for {
		releases, resp, err := client.Repositories.ListReleases(ctx, owner, repo, opt)
		if err != nil {
			return releases, err
		}
		allReleases = append(allReleases, releases...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allReleases, nil
}

// filters all prereleases out of the array (leaves only releases)
func filterPrereleases(releases []*github.RepositoryRelease) []*github.RepositoryRelease {
	var onlyRegularReleases []*github.RepositoryRelease
	for _, release := range releases {
		if !release.GetPrerelease() {
			onlyRegularReleases = append(onlyRegularReleases, release)
		}
	}
	return onlyRegularReleases
}

// filters all releases out of the array (leaves only prereleases)
func filterReleases(releases []*github.RepositoryRelease) []*github.RepositoryRelease {
	var onlyPrereleases []*github.RepositoryRelease
	for _, release := range releases {
		if release.GetPrerelease() {
			onlyPrereleases = append(onlyPrereleases, release)
		}
	}
	return onlyPrereleases
}

// fetches all releases from repo (sorted by publish date)
//
// count specifies the maximum number of releases to return, if its less than 0 there is no max
func getLatestReleases(owner string, repo string, prerelease bool, count int) ([]*github.RepositoryRelease, error) {
	var latestReleases []*github.RepositoryRelease
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
	if count == -1 {
		return latestReleases, nil
	}
	if len(latestReleases) > (count - 1) {
		latestReleases = latestReleases[:count]
	}
	return latestReleases, nil
}
