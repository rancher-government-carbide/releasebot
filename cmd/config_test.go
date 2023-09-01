package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type testMetrics struct {
	testcaseCount   int
	passedTestcases int
	failedTestcases int
}

const repoFile string = "testdata/repos.json"
const payloadFile string = "testdata/payloads.json"

func processResult(t *testing.T, entryField string, passed bool, metrics *testMetrics) {
	if passed {
		metrics.passedTestcases++
	} else {
		t.Errorf("Failed: Incorrect %s field", entryField)
		metrics.failedTestcases++
	}
	metrics.testcaseCount++

}

// checks if the example configuration can be parsed properly, relies on predetermined contents of the config.json file
func Test_loadRepos(t *testing.T) {

	var metrics = testMetrics{
		testcaseCount:   0,
		passedTestcases: 0,
		failedTestcases: 0,
	}
	os.Setenv("RELEASEBOT_REPOS", repoFile)

	var testRepos []RepositoryEntry
	if err := loadRepos(&testRepos); err != nil {
		t.Fatalf("Error loading repo file: %v", err)
	}

	repoIndex := 0
	if testRepos[repoIndex].Owner == "clanktron" {
		processResult(t, "owner", true, &metrics)
	} else {
		processResult(t, "owner", false, &metrics)
	}
	if testRepos[repoIndex].Repo == "dummy" {
		processResult(t, "repo", true, &metrics)
	} else {
		processResult(t, "repo", false, &metrics)
	}
	if testRepos[repoIndex].Slack == true {
		processResult(t, "slack", true, &metrics)
	} else {
		processResult(t, "slack", false, &metrics)
	}
	if testRepos[repoIndex].Prereleases == true {
		processResult(t, "prereleases", true, &metrics)
	} else {
		processResult(t, "prereleases", false, &metrics)
	}
	payloads := []string{"standard", "test"}
	for _, payload := range payloads {
		if !testRepos[repoIndex].Payloads[payload] {
			processResult(t, "payloads", false, &metrics)
		}
	}
	repoIndex++

	if testRepos[repoIndex].Owner == "rancher" {
		processResult(t, "owner", true, &metrics)
	} else {
		processResult(t, "owner", false, &metrics)
	}
	if testRepos[repoIndex].Repo == "rancher" {
		processResult(t, "repo", true, &metrics)
	} else {
		processResult(t, "repo", false, &metrics)
	}
	if testRepos[repoIndex].Slack == true {
		processResult(t, "slack", true, &metrics)
	} else {
		processResult(t, "slack", false, &metrics)
	}
	if testRepos[repoIndex].Prereleases == true {
		processResult(t, "prereleases", true, &metrics)
	} else {
		processResult(t, "prereleases", false, &metrics)
	}
	payloads = []string{"standard", "test"}
	for _, payload := range payloads {
		if !testRepos[repoIndex].Payloads[payload] {
			processResult(t, "payloads", false, &metrics)
		}
	}
	repoIndex++

	if testRepos[repoIndex].Owner == "k3s-io" {
		processResult(t, "owner", true, &metrics)
	} else {
		processResult(t, "owner", false, &metrics)
	}
	if testRepos[repoIndex].Repo == "k3s" {
		processResult(t, "repo", true, &metrics)
	} else {
		processResult(t, "repo", false, &metrics)
	}
	if testRepos[repoIndex].Slack == true {
		processResult(t, "slack", true, &metrics)
	} else {
		processResult(t, "slack", false, &metrics)
	}
	if testRepos[repoIndex].Prereleases == true {
		processResult(t, "prereleases", true, &metrics)
	} else {
		processResult(t, "prereleases", false, &metrics)
	}
	payloads = []string{"standard"}
	for _, payload := range payloads {
		if !testRepos[repoIndex].Payloads[payload] {
			processResult(t, "payloads", false, &metrics)
		}
	}
	repoIndex++

	if testRepos[repoIndex].Owner == "kubernetes" {
		processResult(t, "owner", true, &metrics)
	} else {
		processResult(t, "owner", false, &metrics)
	}
	if testRepos[repoIndex].Repo == "kubernetes" {
		processResult(t, "repo", true, &metrics)
	} else {
		processResult(t, "repo", false, &metrics)
	}
	if testRepos[repoIndex].Slack == true {
		processResult(t, "slack", true, &metrics)
	} else {
		processResult(t, "slack", false, &metrics)
	}
	repoIndex++

	if metrics.failedTestcases > 0 {
		t.Errorf("Failed %d/%d testcases", metrics.failedTestcases, metrics.testcaseCount)
	}

}

func Test_loadPayloads(t *testing.T) {

	correctPayloads := []PayloadEntry{
		{
			Name: "standard",
			Url:  "https://el-example-listener.tekton.svc.cluster.local:8080",
			Payload: json.RawMessage(`
				{
					"Product": "$REPO",
            		"Release": "$RELEASE.TAGNAME"
				}
        `),
		},
		{
			Name: "helm-chart",
			Url:  "https://el-example1-listener.tekton.svc.cluster.local:8080",
			Payload: json.RawMessage(`
				{
					"helm_repo": "$REPO_URL",
            		"release_tag": "$RELEASE.TAGNAME"
				}

			`),
		},
		{
			Name: "example",
			Url:  "https://el-example2-listener.tekton.svc.cluster.local:8080",
			Payload: json.RawMessage(`
				{
					"Other-stuff": "$RELEASE.TAGNAME",
            		"Something-else": "$RELEASE.PUBLISHEDAT",
            		"More_stuff": "AUTHOR.LOGIN"
				}

			`),
		},
	}

	os.Setenv("RELEASEBOT_PAYLOADS", payloadFile)

	var testPayloads []PayloadEntry
	if err := loadPayloads(&testPayloads); err != nil {
		t.Fatalf("Error loading payloads file: %v", err)
	}

	if len(correctPayloads) != len(testPayloads) {
		t.Errorf("Payload arrays have different lengths")
	}

	for i := 0; i < len(correctPayloads) && i < len(testPayloads); i++ {
		if correctPayloads[i].Name != testPayloads[i].Name {
			t.Errorf("Payload name mismatch at index %d: expected %s, got %s", i, correctPayloads[i].Name, testPayloads[i].Name)
		}

		if correctPayloads[i].Url != testPayloads[i].Url {
			t.Errorf("Payload URL mismatch at index %d: expected %s, got %s", i, correctPayloads[i].Url, testPayloads[i].Url)
		}

		require.JSONEq(t, string(correctPayloads[i].Payload), string(testPayloads[i].Payload))
	}

}
