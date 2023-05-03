package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const tektonport int = 8080

func triggertekton(release Release, owner string, repo string) error {

	var tektonurl string = fmt.Sprintf("http://el-%s-release-listener.ssf.svc.cluster.local:%d", repo, tektonport)

	var jsonData = []byte(`{
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
