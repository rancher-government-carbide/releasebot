package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var tekton_namespace = os.Getenv("tekton_namespace")
var tekton_listener = os.Getenv("tekton_listener")
var tekton_port = os.Getenv("tekton_port")

func triggertekton(release Release, owner string, repo string) error {
	var tektonurl string = fmt.Sprintf("http://el-%s-listener.%s.svc.cluster.local:%s", tekton_listener, tekton_namespace, tekton_port)

	var jsonData = []byte(`{
		"product": "` + repo + `",
		"release": "` + release.TagName + `"
	}`)

	req, err := http.NewRequest("POST", tektonurl, bytes.NewBuffer(jsonData))
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
