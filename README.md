# Releasebot
[![Build](https://github.com/clanktron/releasebot/actions/workflows/build.yaml/badge.svg)](https://github.com/clanktron/releasebot/actions/workflows/build.yaml)
[![Unit Tests](https://github.com/clanktron/releasebot/actions/workflows/test.yaml/badge.svg)](https://github.com/clanktron/releasebot/actions/workflows/test.yaml)

A rudimentary daemon that monitors github repos for new releases. 

This is meant to push release events to external sources, currently only slack notifications and tekton pipelines are supported.

Usually this type of event can be pushed by a github action, however if you wish to monitor repos that you don't control then this may come in handy.

## Configuration

| Environment Variable  | Description                                       | Optional          |
| --------------------  | -----------                                       | --------          |
| slack_token           | Oauth token for Your Workspace                    | false             |
| releases_channel      | Channel ID to receive release notifications       | false             |
| prereleases_channel   | Channel ID to receive prerelease notifications    | false             |
| GITHUB_TOKEN          | Github token for authorizing requests             | true              |
| RELEASEBOT_REPOS      | Path to json repo config file                     | true              |
| RELEASEBOT_PAYLOADS   | Path to json payload config file                  | true              |
| interval              | Frequency to query the github api                 | true              |

If the `RELEASEBOT_REPOS` variable is not specified releasebot will read the repos.json in the current directory. It should contain a json array of github repos that you want to monitor.
The format for such is shown below:
```json
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
If the `RELEASEBOT_PAYLOADS` variable is not specified releasebot will read the payloads.json in the current directory. It should contain a json array of github repos that you want to monitor.
The format for such is shown below:
```json
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
            "release_tag": "$RELEASE.TAGNAME"
        }
    },
    {
        "name": "example",
        "url": "https://el-example2-listener.tekton.svc.cluster.local:8080",
        "payload": {
            "Other-stuff": "$RELEASE.TAGNAME",
            "Something-else": "$RELEASE.PUBLISHEDAT",
            "More_stuff": "AUTHOR.LOGIN"
        }
    }
]
```

## Build Binary
```bash
make
```
## Cleanup
```bash
make clean
```
## Help
```bash
make help
```
## Development

#### Git Hooks
```bash
# copy hooks to .git/hooks/
./bin/install-hooks
# symlink hooks to .git/hooks/
./bin/install-hooks link
```
