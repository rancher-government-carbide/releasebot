package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const slackurl string = "https://slack.com/api"

var token = os.Getenv("slack_token")
var releases_channel = os.Getenv("releases_channel")
var prerelease_channel = os.Getenv("prerelease_channel")

func slacknotif(release Release, owner string, repo string, channel string) error {

	if token == "" {
		log.Fatal("Missing slack token")
	}

	publishedDate, err := time.Parse(time.RFC3339, release.PublishedAt)
	if err != nil {
		log.Print(err)
	}

	var jsonData = []byte(`{
		"channel": "` + channel + `",
		"blocks": [
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": "` + owner + `/` + repo + ` -  New Release!"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*` + repo + ` is now at ver.* ` + release.Name + `!\n\n<https://github.com/` + owner + `/` + repo + `/releases/tag/` + release.TagName + `>"
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
				"text": "Authored by: ` + release.Author.Login + ` on ` + publishedDate.Format("Jan 2, 2006") + ` at ` + publishedDate.Format("3:4pm") + `"
			}
			]
		}
		]
	}`)

	req, err := http.NewRequest("POST", slackurl+"/chat.postMessage", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}

// type message struct {
// 	Channel string  `json:"channel"`
// 	Blocks  []block `json:"blocks"`
// }
//
// type block struct {
// 	Type string `json:"type"`
// 	Text slacktext   `json:"text"`
// }
//
// type slackfields struct {
//
// }
//
// type slacktext struct {
// 	Type string `json:"type"`
// 	Text string `json:"text"`
// }
//

//	blocks := []block{
//		{
//			Type: "header",
//			Text: slacktext {
//				Type: "plain_text",
//				Text: "Release TagName",
//			},
//		},
//		{
//			Type: "divider",
//		},
//		{
//			Type: "section",
//			Fields:
//		},
//	}
//
//	completemessage := message{
//		Channel: slackchannel,
//		Blocks: blocks,
//	}o
