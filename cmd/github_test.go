package main

import (
	"testing"
	"time"

	"github.com/google/go-github/v55/github"
)

var testRelease1 = &github.RepositoryRelease{
	Name:        github.String("v1.0.0"),
	TagName:     github.String("v1.0.0"),
	Body:        github.String("This is the release description"),
	Draft:       github.Bool(false),
	Prerelease:  github.Bool(false),
	PublishedAt: &github.Timestamp{Time: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)},
}

var testRelease2 = &github.RepositoryRelease{
	Name:        github.String("v1.1.0"),
	TagName:     github.String("v1.1.0"),
	Body:        github.String("This is release v1.1.0"),
	Draft:       github.Bool(false),
	Prerelease:  github.Bool(false),
	PublishedAt: &github.Timestamp{Time: time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC)},
}

var testRelease3 = &github.RepositoryRelease{
	Name:        github.String("v2.0.0"),
	TagName:     github.String("v2.0.0"),
	Body:        github.String("This is release v2.0.0"),
	Draft:       github.Bool(false),
	Prerelease:  github.Bool(false),
	PublishedAt: &github.Timestamp{Time: time.Date(2023, 3, 25, 0, 0, 0, 0, time.UTC)},
}

var testRelease4 = &github.RepositoryRelease{
	Name:        github.String("v2.1.0"),
	TagName:     github.String("v2.1.0"),
	Body:        github.String("This is release v2.1.0"),
	Draft:       github.Bool(false),
	Prerelease:  github.Bool(false),
	PublishedAt: &github.Timestamp{Time: time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC)},
}

var testRelease5 = &github.RepositoryRelease{
	Name:        github.String("v3.0.0"),
	TagName:     github.String("v3.0.0"),
	Body:        github.String("This is release v3.0.0"),
	Draft:       github.Bool(false),
	Prerelease:  github.Bool(false),
	PublishedAt: &github.Timestamp{Time: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)},
}

// const testReleaseCount int = 5

var testReleases = []*github.RepositoryRelease{
	testRelease1,
	testRelease2,
	testRelease3,
	testRelease4,
	testRelease5,
}

func TestSortByPublishedDate(t *testing.T) {
	expectedTimestamps := []time.Time{
		time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), // v3.0.0
		time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC), // v2.1.0
		time.Date(2023, 3, 25, 0, 0, 0, 0, time.UTC), // v2.0.0
		time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC), // v1.1.0
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), // v1.0.0
	}

	sortedReleases := sortByPublishDate(testReleases)

	for i, release := range sortedReleases {
		if release.PublishedAt == nil {
			t.Errorf("Release %d has a nil PublishedAt field.", i+1)
			continue
		}
		if release.GetPublishedAt().Time != expectedTimestamps[i] {
			t.Errorf("Expected release %d to have PublishedAt timestamp %v, but got %v.", i+1, expectedTimestamps[i], release.GetPublishedAt().String())
		}
	}
}

// func TestGetReleases(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Check the request URL
// 		expectedURL := fmt.Sprintf("/repos/%s/%s/releases", "owner", "repo")
// 		if r.URL.Path != expectedURL {
// 			t.Errorf("Expected URL path to be %s, got %s", expectedURL, r.URL.Path)
// 			return
// 		}
// 		// Check the authorization header
// 		expectedToken := "test_token"
// 		token := r.Header.Get("Authorization")
// 		if token != fmt.Sprintf("token %s", expectedToken) {
// 			t.Errorf("Expected Authorization header to be 'token %s', got '%s'", expectedToken, token)
// 			return
// 		}
// 		// Return a sample response
// 		err := json.NewEncoder(w).Encode(testReleases)
// 		if err != nil {
// 			t.Errorf("Failed to encode response: %s", err)
// 			return
// 		}
// 	}))
// 	defer ts.Close()
//
// 	github_api_url = ts.URL
// 	os.Setenv("GITHUB_TOKEN", "test_token")
// 	releases, err := getAllReleases("owner", "repo")
// 	if err != nil {
// 		t.Errorf("Failed to get releases: %s", err)
// 	}
// 	expectedCount := test_release_count
// 	if len(releases) != expectedCount {
// 		t.Errorf("Expected %d releases, got %d", expectedCount, len(releases))
// 	}
// 	if !reflect.DeepEqual(releases, testReleases) {
// 		t.Errorf("Expected releases to be %v, got %v", testReleases, releases)
// 	}
// }

func TestFilterPrereleases(t *testing.T) {
	filteredReleases := filterPrereleases(testReleases)
	for _, release := range filteredReleases {
		if release.GetPrerelease() {
			t.Errorf("FilterPrereleases failed - found prerelease in filtered releases: %s", release.GetTagName())
		}
	}
}

func TestFilterReleases(t *testing.T) {
	filteredReleases := filterReleases(testReleases)
	for _, release := range filteredReleases {
		if !release.GetPrerelease() {
			t.Errorf("FilterReleases failed - found release in filtered releases: %s", release.GetTagName())
		}
	}
}
