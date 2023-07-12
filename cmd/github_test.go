package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

var test_release_1 = Release{
	Name:        "v2.3.4",
	TagName:     "v2.3.4",
	Prerelease:  false,
	HtmlUrl:     "https://example.com/releases/v2.3.4",
	PublishedAt: "2023-01-10T12:34:56Z",
	Author: Author{
		Login:     "john_doe",
		AvatarUrl: "https://example.com/avatar.jpg",
		HtmlUrl:   "https://example.com/users/john_doe",
	},
}

var test_release_2 = Release{
	Name:        "v2.3.4",
	TagName:     "v2.3.4",
	Prerelease:  true,
	HtmlUrl:     "https://example.com/releases/v2.3.4-beta",
	PublishedAt: "2023-05-21T08:15:30Z",
	Author: Author{
		Login:     "jane_smith",
		AvatarUrl: "https://example.com/avatar.jpg",
		HtmlUrl:   "https://example.com/users/jane_smith",
	},
}

var test_release_3 = Release{
	Name:        "v2.3.4",
	TagName:     "v2.3.4",
	Prerelease:  false,
	HtmlUrl:     "https://example.com/releases/v2.3.4",
	PublishedAt: "2024-03-12T10:20:30Z",
	Author: Author{
		Login:     "alice_walker",
		AvatarUrl: "https://example.com/avatar.jpg",
		HtmlUrl:   "https://example.com/users/alice_walker",
	},
}

var test_release_4 = Release{
	Name:        "v2.3.4",
	TagName:     "v2.3.4",
	Prerelease:  true,
	HtmlUrl:     "https://example.com/releases/v2.3.4-beta",
	PublishedAt: "2024-08-28T15:45:00Z",
	Author: Author{
		Login:     "bob_jackson",
		AvatarUrl: "https://example.com/avatar.jpg",
		HtmlUrl:   "https://example.com/users/bob_jackson",
	},
}

var test_release_5 = Release{
	Name:        "v2.3.4",
	TagName:     "v2.3.4",
	Prerelease:  false,
	HtmlUrl:     "https://example.com/releases/v2.3.4",
	PublishedAt: "2025-06-17T09:00:15Z",
	Author: Author{
		Login:     "samuel_rodriguez",
		AvatarUrl: "https://example.com/avatar.jpg",
		HtmlUrl:   "https://example.com/users/samuel_rodriguez",
	},
}

var test_releases = []Release{
	test_release_1,
	test_release_2,
	test_release_3,
	test_release_4,
	test_release_5,
}

const test_release_count int = 5

func TestSortByPublishedDate(t *testing.T) {
	expectedOrder := []string{
		"2025-06-17T09:00:15Z",
		"2024-08-28T15:45:00Z",
		"2024-03-12T10:20:30Z",
		"2023-05-21T08:15:30Z",
		"2023-01-10T12:34:56Z",
	}

	sortedReleases := sort_by_published_date(test_releases)

	for i, release := range sortedReleases {
		if release.PublishedAt != expectedOrder[i] {
			t.Errorf("SortByPublishedDate failed - expected: %s, got: %s", expectedOrder[i], release.PublishedAt)
		}
	}
}

func TestGetReleases(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request URL
		expectedURL := fmt.Sprintf("/repos/%s/%s/releases", "owner", "repo")
		if r.URL.Path != expectedURL {
			t.Errorf("Expected URL path to be %s, got %s", expectedURL, r.URL.Path)
			return
		}
		// Check the authorization header
		expectedToken := "test_token"
		token := r.Header.Get("Authorization")
		if token != fmt.Sprintf("token %s", expectedToken) {
			t.Errorf("Expected Authorization header to be 'token %s', got '%s'", expectedToken, token)
			return
		}
		// Return a sample response
		err := json.NewEncoder(w).Encode(test_releases)
		if err != nil {
			t.Errorf("Failed to encode response: %s", err)
			return
		}
	}))
	defer ts.Close()

	github_api_url = ts.URL
	os.Setenv("github_token", "test_token")
	releases, err := get_releases("owner", "repo")
	if err != nil {
		t.Errorf("Failed to get releases: %s", err)
	}
	expectedCount := test_release_count
	if len(releases) != expectedCount {
		t.Errorf("Expected %d releases, got %d", expectedCount, len(releases))
	}
	if !reflect.DeepEqual(releases, test_releases) {
		t.Errorf("Expected releases to be %v, got %v", test_releases, releases)
	}
}

func TestFilterPrereleases(t *testing.T) {
	filteredReleases := filter_prereleases(test_releases)

	for _, release := range filteredReleases {
		if release.Prerelease {
			t.Errorf("FilterPrereleases failed - found prerelease in filtered releases: %s", release.TagName)
		}
	}
}

func TestFilterReleases(t *testing.T) {
	filteredReleases := filter_releases(test_releases)

	for _, release := range filteredReleases {
		if !release.Prerelease {
			t.Errorf("FilterReleases failed - found release in filtered releases: %s", release.TagName)
		}
	}
}
