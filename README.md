# Releasebot
[![Binary Build and Test](https://github.com/clanktron/releasebot/actions/workflows/go-build.yaml/badge.svg)](https://github.com/clanktron/releasebot/actions/workflows/go-build.yaml)

A rudimentary daemon that monitors github repos for new releases. 

This is meant to push release events to external sources, currently only slack notifications and tekton pipelines are supported.

Usually this type of event can be pushed by a github action, however if you wish to monitor repos that you don't control then this may come in handy.

Roadmap (ordered most to least important):
- helm chart

## Configuration

| Environment Variable  | Description                                       | Optional          |
| --------------------  | -----------                                       | --------          |
| slack_token           | Oauth token for Your Workspace                    | false             |
| releases_channel      | Channel ID to receive release notifications       | false             |
| prereleases_channel   | Channel ID to receive prerelease notifications    | false             |
| github_token          | Github token for authorizing requests             | true              |
| releasebot_config     | Path to json config file                          | true              |
| interval              | Frequency to query the github api                 | true              |
| tekton_listener       | If utilizing Tekton, the event listener name      | true              |
| tekton_namespace      | If utilizing Tekton, the namespace of EL          | true              |
| tekton_port           | If utilizing Tekton, the port of the EL           | true              |

If the `releasebot_config` variable is not specified releasebot will read the config.json in the current directory. It should contain a json array of github repos that you want to monitor.
The format for such is shown below:
```json
[
    {
        "owner": "owner",
        "repo": "repo",
        "prereleases": true,
        "tekton": true
    },
    {
        "owner": "clanktron",
        "repo": "releasebot"
        "prereleases": true,
        "tekton": true
    },
    {
        "owner": "golang",
        "repo": "go"
        "prereleases": false,
        "tekton": false
    }
]
```

## Build Binary
```bash
make
```
## Test Binary
```bash
make test
```
## Build Container
```bash
make container
```
## Cleanup
```bash
make clean
```
## Help
```bash
make help
```
