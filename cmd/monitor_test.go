package main

import (
	"reflect"
	"sort"
	"testing"
)

func Test_stringifyLoadedReleases(t *testing.T) {

	testloadedReleasesMap := map[string]bool{
		"v3.4.5":        true,
		"release-2.0":   true,
		"v324":          true,
		"v893":          true,
		"ungabunga":     true,
		"notinthearray": false,
	}

	correctStringifiedReleases := []string{
		"v3.4.5",
		"release-2.0",
		"v324",
		"v893",
		"ungabunga",
	}

	testStringifiedReleases := stringifyLoadedReleases(testloadedReleasesMap)

	sort.Strings(testStringifiedReleases)
	sort.Strings(correctStringifiedReleases)

	if !reflect.DeepEqual(correctStringifiedReleases, testStringifiedReleases) {
		t.Errorf("Arrays of release names do not match")
		t.Logf("Correct array: %v", correctStringifiedReleases)
		t.Logf("Test array: %v", testStringifiedReleases)
	}

}
