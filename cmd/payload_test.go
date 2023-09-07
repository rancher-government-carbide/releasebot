package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/google/go-github/v55/github"
)

func TestReplaceVariables(t *testing.T) {
	tests := []struct {
		name      string
		data      map[string]interface{}
		variables map[string]string
		expected  map[string]interface{}
	}{
		{
			name: "ReplaceVariables - Basic",
			data: map[string]interface{}{
				"name": "$firstName",
				"age":  "$age",
			},
			variables: map[string]string{
				"firstName": "John",
				"age":       "30",
			},
			expected: map[string]interface{}{
				"name": "John",
				"age":  "30",
			},
		},
		{
			name: "ReplaceVariables - Nested Map",
			data: map[string]interface{}{
				"info": map[string]interface{}{
					"name": "$firstName",
					"age":  "$age",
				},
			},
			variables: map[string]string{
				"firstName": "Alice",
				"age":       "25",
			},
			expected: map[string]interface{}{
				"info": map[string]interface{}{
					"name": "Alice",
					"age":  "25",
				},
			},
		},
		{
			name: "ReplaceVariables - No Variables",
			data: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			variables: map[string]string{
				"firstName": "Alice",
				"city":      "Paris",
			},
			expected: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			replaceVariables(test.data, test.variables)
			if !reflect.DeepEqual(test.data, test.expected) {
				t.Errorf("ReplaceVariables did not produce the expected result for test case '%s'. Got: %+v, Expected: %+v", test.name, test.data, test.expected)
			}
		})
	}
}

func sortMap(m map[string]interface{}) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := m[k]
		if subMap, ok := v.(map[string]interface{}); ok {
			sortMap(subMap)
		}
	}
}

func TestParsePayload(t *testing.T) {

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
            		"RepoUrl": "$REPO.URL",
            		"Tagname": "$RELEASE.TAGNAME",
            		"Prerelease": "$RELEASE.PRERELEASE",
            		"HtmlUrl": "$RELEASE.HTMLURL",
            		"PublishedAt": "$RELEASE.PUBLISHEDAT",
            		"Author.Login": "$AUTHOR.LOGIN",
            		"Author.AvatarUrl": "$AUTHOR.AVATARURL",
            		"Author.HtmlUrl": "$AUTHOR.HTMLURL"
				}
			`),
	}

	var testRelease = &github.RepositoryRelease{
		Name:        github.String("v2.3.4"),
		TagName:     github.String("v2.3.4"),
		Prerelease:  github.Bool(false),
		HTMLURL:     github.String("https://example.com/releases/v2.3.4"),
		PublishedAt: &github.Timestamp{Time: time.Date(2023, 7, 16, 12, 0, 0, 0, time.UTC)},
		Author: &github.User{
			Login:     github.String("john_doe"),
			AvatarURL: github.String("https://example.com/avatar.jpg"),
			HTMLURL:   github.String("https://example.com/users/john_doe"),
		},
	}

	correctCompiledPayload := json.RawMessage(`
				{
					"Product": "example_repo",
					"RepoUrl": "git@github.com:example_owner/example_repo",
            		"Tagname": "v2.3.4",
            		"Prerelease": "false",
            		"HtmlUrl": "https://example.com/releases/v2.3.4",
            		"PublishedAt": "2023-07-16 12:00:00 +0000 UTC",
            		"Author.Login": "john_doe",
            		"Author.AvatarUrl": "https://example.com/avatar.jpg",
            		"Author.HtmlUrl": "https://example.com/users/john_doe"
				}
			`)

	testCompiledJSONPayload, err := parsePayload(testRelease, testRepoEntry, testPayloadEntry)
	if err != nil {
		t.Fatalf("Failed to parse payloads: %v", err)
	}

	// Unmarshal the testCompiledJSONPayload
	var testPayload map[string]interface{}
	err = json.Unmarshal(testCompiledJSONPayload, &testPayload)
	if err != nil {
		t.Fatalf("Failed to unmarshal testCompiledJSONPayload: %v", err)
	}

	// Unmarshal the correctCompiledJSONPayload
	var correctPayload map[string]interface{}
	err = json.Unmarshal(correctCompiledPayload, &correctPayload)
	if err != nil {
		t.Fatalf("Failed to unmarshal correctCompiledJSONPayload: %v", err)
	}

	// Sort the elements within the maps
	sortMap(testPayload)
	sortMap(correctPayload)

	// Marshal the sorted maps back into JSON
	testJSON, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("Failed to marshal testPayload: %v", err)
	}

	correctJSON, err := json.Marshal(correctPayload)
	if err != nil {
		t.Fatalf("Failed to marshal correctPayload: %v", err)
	}

	// Compare the JSON strings
	if string(testJSON) != string(correctJSON) {
		t.Error("JSON payload compiled improperly")
		t.Logf("Expected JSON Payload: %s", correctJSON)
		t.Logf("Actual JSON Payload: %s", testJSON)
	}

}

func TestSendPayload(t *testing.T) {
	// Create a test server to mock the HTTP request/response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading request body: %v", err)
		}

		// Check the request body
		expectedPayload := []byte(`{"key": "value"}`)
		if !bytes.Equal(body, expectedPayload) {
			t.Errorf("Expected request body %s, got %s", expectedPayload, body)
		}

		// Send a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	// Call the sendPayload function with the test server URL
	err := sendPayload([]byte(`{"key": "value"}`), server.URL)
	if err != nil {
		t.Fatalf("sendPayload failed: %v", err)
	}
}
