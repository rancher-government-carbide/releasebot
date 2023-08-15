package main

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func Test_parsePayload(t *testing.T) {

	correctCompiledPayload := json.RawMessage(`
				{
					"Product": "example_repo",
            		"Release": "v2.3.4"
				}
			`)

	testRepoEntry := RepositoryEntry{
		Owner: "example_owner",
		Repo:  "example_repo",
	}

	testPayloadEntry := PayloadEntry{
		Name: "example_payload",
		Url:  "https://example.com/payload",
		Payload: json.RawMessage(`
				{
					"Product": "$REPO",
            		"Release": "$RELEASE.TAGNAME"
				}
			`),
	}

	var testRelease = Release{
		Name:       "v2.3.4",
		TagName:    "v2.3.4",
		Prerelease: false,
		HtmlUrl:    "https://example.com/releases/v2.3.4",
		PublishedAt: Time{
			Time: time.Date(2023, 7, 16, 12, 0, 0, 0, time.UTC),
		},
		Author: Author{
			Login:     "john_doe",
			AvatarUrl: "https://example.com/avatar.jpg",
			HtmlUrl:   "https://example.com/users/john_doe",
		},
	}

	testCompiledJSONPayload, err := parsePayload(testRelease, testRepoEntry, testPayloadEntry)
	if err != nil {
		t.Fatalf("Failed to parse payloads: %v", err)
	}
	correctCompiledJSONPayload, err := json.Marshal(correctCompiledPayload)
	if err != nil {
		t.Fatalf("Failed to marshall correctCompiledPayload: %v\n Test invalid", err)
	}

	if !reflect.DeepEqual(testCompiledJSONPayload, correctCompiledJSONPayload) {
		t.Error("JSON payload compiled improperly")
		t.Logf("Expected JSON Payload: %s", correctCompiledJSONPayload)
		t.Logf("Actual JSON Payload: %s", testCompiledJSONPayload)
	}

}
