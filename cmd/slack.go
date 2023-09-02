package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var slackurl string = "https://slack.com/api"

var token = os.Getenv("slack_token")
var releases_channel = os.Getenv("releases_channel")
var prereleases_channel = os.Getenv("prereleases_channel")

func slacknotif(release Release, owner string, repo string) error {

	if token == "" {
		log.Fatal("Missing slack token")
	}

	publishedDate := release.PublishedAt.Time
	release_type := "Release"
	channel := releases_channel
	if release.Prerelease {
		release_type = "Prerelease"
		channel = prereleases_channel
	}

	var jsonData = []byte(`{
		"channel": "` + channel + `",
		"blocks": [
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": "` + owner + `/` + repo + ` -  New ` + release_type + `!"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "` + release.Name + ` is now available!\n\n<https://github.com/` + owner + `/` + repo + `/releases/tag/` + release.TagName + `>"
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
				"image_url": "` + release.Author.AvatarUrl + `",
				"alt_text": "author profile img"
			},
			{
				"type": "mrkdwn",
				"text": "Authored by: ` + release.Author.Login + ` on ` + publishedDate.Format("Jan 2, 2006") + ` at ` + publishedDate.In(time.UTC).Format("3:4pm MST") + `"
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
	body, _ := io.ReadAll(resp.Body)
	log.Println("slack response Status:", resp.Status)
	log.Println("slack response Body:", string(body))

	return nil
}
