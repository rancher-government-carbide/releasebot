package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v55/github"
)

var slackurl string = "https://slack.com/api"

var token = os.Getenv("slack_token")
var releasesChannel = os.Getenv("releases_channel")
var prereleasesChannel = os.Getenv("prereleases_channel")

func slacknotif(release *github.RepositoryRelease, owner string, repo string) error {

	if token == "" {
		log.Fatal("Missing slack token")
	}

	publishedDate := release.GetPublishedAt().Time
	releaseType := "Release"
	channel := releasesChannel
	if release.GetPrerelease() {
		releaseType = "Prerelease"
		channel = prereleasesChannel
	}

	var jsonData = []byte(`{
		"channel": "` + channel + `",
		"blocks": [
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": "` + owner + `/` + repo + ` -  New ` + releaseType + `!"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "` + release.GetName() + ` is now available!\n\n<https://github.com/` + owner + `/` + repo + `/releases/tag/` + release.GetTagName() + `>"
			},
			"accessory": {
				"type": "image",
				"image_url": "https://github.com/` + owner + `.png",
				"alt_text": "repo icon"
			}
		},
		{
			"type": "context",
			"elements": [
			{
				"type": "image",
				"image_url": "` + release.Author.GetAvatarURL() + `",
				"alt_text": "author profile img"
			},
			{
				"type": "mrkdwn",
				"text": "Authored by: ` + release.Author.GetLogin() + ` on ` + publishedDate.Format("Jan 2, 2006") + ` at ` + publishedDate.In(time.UTC).Format("3:04pm MST") + `"
			}
			]
		}
		]
	}`)

	req, err := http.NewRequest("POST", slackurl+"/chat.postMessage", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error sending slack message - response status: %s - response body : %s", resp.Status, string(body))
	}
	return nil
}
