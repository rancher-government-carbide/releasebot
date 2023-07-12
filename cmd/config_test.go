package main

import (
	"os"
	"testing"
)

const configfile string = "../config.json"

// checks if the example configuration can be parsed properly, relies on predetermined contents of the config.json file
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
	var entry string
	const passed_testcase string = "Passed %s testcase"
	const failed_testcase string = "Failed %s testcase"

	entry = "clanktron/dummy"
	if config[testcase].Owner == "clanktron" && config[testcase].Repo == "dummy" && config[testcase].Prereleases && config[testcase].Tekton {
		t.Logf(passed_testcase, entry)
		passed++
	} else {
		t.Logf(failed_testcase, entry)
		failed++
	}
	testcase++

	entry = "rancher/rancher"
	if config[testcase].Owner == "rancher" && config[testcase].Repo == "rancher" {
		t.Logf(passed_testcase, entry)
		passed++
	} else {
		t.Logf(failed_testcase, entry)
		failed++
	}
	testcase++

	if config[testcase].Owner == "k3s-io" && config[testcase].Repo == "k3s" {
		t.Logf(passed_testcase, entry)
		passed++
	} else {
		t.Logf(failed_testcase, entry)
		failed++
	}
	testcase++

	if config[testcase].Owner == "kubernetes" && config[testcase].Repo == "kubernetes" {
		t.Logf(passed_testcase, entry)
		passed++
	} else {
		t.Logf(failed_testcase, entry)
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
