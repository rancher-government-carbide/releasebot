package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const slackurl string = "https://slack.com/api"

func slacknotif(release Release) error {

	token := os.Getenv("slack_token")
	if token == "" {
		log.Fatal("Missing slack token")
	}

	channel := os.Getenv("slack_channel")
	if channel == "" {
		log.Fatal("Missing slack channel ID")
	}

	var jsonData = []byte(`{
		"channel": "` + channel + `",
		"blocks": [
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": "` + release.TagName + `"
				}
			},
			{
				"type": "divider"
			},
			{
				"type": "section",
				"fields": [
					{
						"type": "mrkdwn",
						"text": "*Current Quarter*\nBudget: $18,000 (ends in 53 days)\nSpend: $4,289.70\nRemain: $13,710.30"
					},
					{
						"type": "mrkdwn",
						"text": "*Top Expense Categories*\n:airplane: Flights · 30%\n:taxi: Taxi / Uber / Lyft · 24% \n:knife_fork_plate: Client lunch / meetings · 18%"
					}
				]
			},
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
