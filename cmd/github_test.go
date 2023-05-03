package main

import (
	"testing"
)

var test_release_1 = Release{
	Name:       "test1",
	TagName:    "test1",
	Prerelease: false,
}

var test_release_2 = Release{
	Name:       "test2",
	TagName:    "test2",
	Prerelease: false,
}

var test_release_3 = Release{
	Name:       "test3",
	TagName:    "test3",
	Prerelease: false,
}

var test_releases = []Release{test_release_1, test_release_2, test_release_3}

func Test_get_releases(t *testing.T) {
	// basic get request - I ain't writing a test for that
}

func Test_filter_prerelease(t *testing.T) {

	// var filtered_releases []Release
	// filtered_releases = filter_prereleases(test_releases)

	// if filtered_releases != []Release{test_release_3} {

	// }

}

func Test_get_latest_release(t *testing.T) {

}

func Test_releasecheck(t *testing.T) {

}
