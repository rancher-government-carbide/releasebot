package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var tekton_namespace = os.Getenv("tekton_namespace")
var tekton_listener = os.Getenv("tekton_listener")
var tekton_port = os.Getenv("tekton_port")

type publish_chart_trigger struct {
	Helm_repo   string `json:"helm_repo"`
	Release_tag string `json:"release_tag"`
}

func triggertekton(release Release, owner string, repo string) error {

	if tekton_namespace == "" {
		tekton_namespace = "tekton-pipelines"
	}
	if tekton_listener == "" {
		tekton_listener = repo
	}
	if tekton_port == "" {
		tekton_port = "8080"
	}

	var tektonurl string = fmt.Sprintf("http://el-%s-listener.%s.svc.cluster.local:%s", tekton_listener, tekton_namespace, tekton_port)
	var repo_url string = fmt.Sprintf("git@github.com:%s/%s", owner, repo)

	trigger := publish_chart_trigger{
		Helm_repo:   repo_url,
		Release_tag: release.TagName,
	}

	json_trigger, err := json.Marshal(trigger)

	req, err := http.NewRequest("POST", tektonurl, bytes.NewBuffer(json_trigger))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return err
	} else {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Body:", string(body))
	}

	return nil
}
