# Releasebot
[![Build](https://github.com/clanktron/releasebot/actions/workflows/build.yaml/badge.svg)](https://github.com/clanktron/releasebot/actions/workflows/build.yaml)
[![Unit Tests](https://github.com/clanktron/releasebot/actions/workflows/test.yaml/badge.svg)](https://github.com/clanktron/releasebot/actions/workflows/test.yaml)

A rudimentary daemon that monitors github repos for new releases. 

This is meant to push release events to external sources.
Usually this type of event can be pushed by a github action, however if you wish to monitor repos that you don't control then this may come in handy.

## Configuration

### Environment
| Environment Variable  | Description                                                                       | Optional  |
| --------------------  | -----------                                                                       | --------  |
| slack_token           | Oauth token for Your Workspace                                                    | false     |
| releases_channel      | Channel ID to receive release notifications                                       | false     |
| prereleases_channel   | Channel ID to receive prerelease notifications                                    | false     |
| GITHUB_TOKEN          | Github token for authorizing requests                                             | true      |
| RELEASEBOT_REPOS      | Path to json repo config file                                                     | true      |
| RELEASEBOT_PAYLOADS   | Path to json payload config file                                                  | true      |
| PERSIST               | Set to "true" or "TRUE" if you wish to track releases across releasebot restarts  | true      |
| interval              | Frequency to query the github api                                                 | true      |

### Config Files
If the `RELEASEBOT_REPOS` variable is not specified releasebot will read the repos.json in the current directory.
It should contain a json array of github repos that you want to monitor.
The members of the payloads array should correspond to entries in the payloads.json file.
The format for such is shown below:

repos.json:
```repos.json
[
    {
        "owner": "clanktron",
        "repo": "dummy",
        "slack": true,
        "prereleases": true,
        "payloads": [ "standard", "helm-chart" ]
    },
    {
        "owner": "rancher",
        "repo": "rancher",
        "slack": true,
        "prereleases": true,
        "payloads": [ "standard", "example" ]
    },
    {
        "owner": "k3s-io",
        "repo": "k3s",
        "slack": true,
        "prereleases": true,
        "payloads": [ "standard" ]
    },
    {
        "owner": "kubernetes",
        "repo": "kubernetes",
        "slack": true
    }
]
```
#### Fields:
- **owner (string):** The owner or organization name of the GitHub repository.
- **repo (string):** The name of the GitHub repository.
- **slack (boolean, optional):** A flag indicating whether Slack notifications are enabled for this repository. It can be true or false (defaults to false).
- **prereleases (boolean, optional):** A flag indicating whether pre-releases should be monitored as well for this repository. It can be true or false (defaults to false).
- **payloads (array of strings):** an array of payload types associated with this repository. Possible values include any names of payloads specified in payloads.json.

If the `RELEASEBOT_PAYLOADS` variable is not specified releasebot will read the payloads.json in the current directory.
It should contain an array of the json payloads you wish to send to the specified urls (webhooks you wish to trigger etc).
The format for such is shown below:

payloads.json:
```payloads.json
[
    {
        "name": "standard",
        "url": "https://el-example-listener.tekton.svc.cluster.local:8080",
        "payload": {
            "Product": "$REPO",
            "Release": "$RELEASE.TAGNAME"
        }
    },
    {
        "name": "helm-chart",
        "url": "https://el-example1-listener.tekton.svc.cluster.local:8080",
        "payload": {
            "helm_repo": "$REPO_URL",
            "release_tag": "$RELEASE.TAGNAME",
            "otherData": "something hardcoded"
        }
    },
    {
        "name": "example",
        "url": "https://el-example2-listener.tekton.svc.cluster.local:8080",
        "payload": {
            "Other-stuff": "$RELEASE.TAGNAME",
            "Something-else": "$RELEASE.PUBLISHEDAT",
            "More_stuff": "$AUTHOR.LOGIN"
        }
    }
]
```
#### Fields:
- **name (string):** The name of the json payload to be referenced in repos.json.
- **url (string):** The url you wish to send your json to.
- **payload (json object):** A JSON object that you want sent to the address specified in the url field. It can be any valid json. 
Certain variables are available for runtime substitution if you need information about the release in your json payload. 
These must be all caps and be prefixed with a `$`.

Available Variables:

| Variable              | Description
| --------------------  | -----------
| $REPO                  | Name of the repository
| $REPO.URL              | ssh url of the repository
| $RELEASE.TAGNAME       | Tag corresponding to the release
| $RELEASE.PRERELEASE    | Stringified boolean of whether release is a prerelease
| $RELEASE.HTMLURL       | Url for viewing the release on Github
| $RELEASE.PUBLISHEDAT   | Date+Time the release was published at
| $AUTHOR.LOGIN          | Username of the release author
| $AUTHOR.AVATARURL      | Url for viewing the Github avatar image of the release author
| $AUTHOR.HTMLURL        | Url for viewing the Github account of the release author

## Helm

Add public source:
```bash
helm repo add releasebot https://rancher-government-carbide.github.io/releasebot
helm repo update
```

Install:
```bash
export HELM_RELEASE_NAME=releasebot
export VALUES_FILE=values.yaml
# from public release
helm install $HELM_RELEASE_NAME releasebot/releasebot --values $VALUES_FILE
# OR
# with locally cloned repo
helm install $HELM_RELEASE_NAME ./chart --values $VALUES_FILE
```

## Development

#### Build
```bash
make help
# Available targets:
#   releasebot            Build the binary (default)
#   test                  Run all unit tests
#   lint                  Run go vet and staticcheck
#   check                 Build, test, and lint the binary
#   linux                 Build the binary for Linux
#   darwin                Build the binary for MacOS
#   windows               Build the binary for Windows
#   container             Build the container
#   container-push        Build and push the container
#   clean                 Clean build results
#   help                  Show help
```

#### Git Hooks
```bash
# copy hooks to .git/hooks/
./bin/install-hooks
# symlink hooks to .git/hooks/
./bin/install-hooks link
```
