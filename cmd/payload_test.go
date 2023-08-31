package main

import (
	"encoding/json"
	"reflect"
	"bytes"
	"testing"
	"time"
	"io"
	"net/http"
	"net/http/httptest"
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
