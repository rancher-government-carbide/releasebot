package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v55/github"
)

func replaceVariables(data map[string]interface{}, variables map[string]string) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			// Check if the string contains a variable
			if strings.HasPrefix(v, "$") {
				varName := v[1:] // remove the $
				if val, ok := variables[varName]; ok {
					data[key] = val
				}
			}
		case map[string]interface{}:
			// Recursively process nested maps
			replaceVariables(v, variables)
		}
	}
}

func parsePayload(release *github.RepositoryRelease, repo RepositoryEntry, payload PayloadEntry) ([]byte, error) {

	var repo_url string = fmt.Sprintf("git@github.com:%s/%s", repo.Owner, repo.Repo)

	variables := map[string]string{
		"REPO":                repo.Repo,
		"REPO.URL":            repo_url,
		"RELEASE.TAGNAME":     release.GetTagName(),
		"RELEASE.PRERELEASE":  strconv.FormatBool(release.GetPrerelease()),
		"RELEASE.HTMLURL":     release.GetHTMLURL(),
		"RELEASE.PUBLISHEDAT": release.GetPublishedAt().String(),
		"AUTHOR.LOGIN":        release.Author.GetLogin(),
		"AUTHOR.AVATARURL":    release.Author.GetAvatarURL(),
		"AUTHOR.HTMLURL":      release.Author.GetHTMLURL(),
	}

	var data map[string]interface{}
	if err := json.Unmarshal(payload.Payload, &data); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to unmarshall payload")
	}

	replaceVariables(data, variables)

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonPayload, nil
}

func sendPayload(jsonPayload []byte, url string) error {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"body":  body,
			"url":   url,
			"error": err,
		}).Error("Failed to deliver payload to url")
		return err
	}

	return nil
}

func sendAllPayloads(release *github.RepositoryRelease, repo RepositoryEntry, payloadEntries []PayloadEntry) error {
	for _, payload := range payloadEntries {
		if repo.Payloads[payload.Name] {
			renderedPayload, err := parsePayload(release, repo, payload)
			if err != nil {
				return err
			}
			err = sendPayload(renderedPayload, payload.Url)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
