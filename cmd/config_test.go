package main

import (
	"fmt"
	"os"
	"testing"
)

const configfile string = "../config.json"

func Test_loadConfig(t *testing.T) {

	os.Setenv("releasebot_config", configfile)

	var config []Repository
	if err := loadConfig(&config); err != nil {
		// t.Logf("Error loading config file...\nExiting...\n")
	} else {
		t.Logf("Loaded config file successfully...\n")
	}

	var testcase int = 0
	var passed int = 0
	var failed int = 0

	if config[testcase].Owner == "clanktron" && config[testcase].Repo == "dummy" && config[testcase].Prereleases && config[testcase].Tekton {
		// t.Logf("Passed Test: %d ", testcase+1)
		passed++
	} else {
		fmt.Println(config[testcase])
		failed++
	}
	testcase++

	if config[testcase].Owner == "rancher" && config[testcase].Repo == "rancher" {
		// t.Logf("Passed Test")
		passed++
	} else {
		// t.Logf("Failed Test")
		failed++
	}
	testcase++

	if config[testcase].Owner == "k3s-io" && config[testcase].Repo == "k3s" {
		// t.Logf("Passed Test")
		passed++
	} else {
		// t.Logf("Failed Test")
		failed++
	}
	testcase++

	if config[testcase].Owner == "kubernetes" && config[testcase].Repo == "kubernetes" {
		// t.Logf("Passed Test")
		passed++
	} else {
		// t.Logf("Failed Test")
		failed++
	}
	testcase++

	if passed > 0 {
		t.Logf("Passed %d/%d testcases", passed, testcase)
	}
	if failed > 0 {
		t.Errorf("Failed %d/%d testcases", failed, testcase)
	}

}
