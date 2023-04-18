# Releasebot

A rudimentary daemon that monitors github repos for new releases. 

This is meant to push release events to external sources, currently only slack notifications are supported.

Triggering Tekton pipelines is on the roadmap, along with any other use cases that present themselves.

Usually this type of event can be pushed by a github action, however if you wish to monitor repos that you don't control then this may come in handy.

Roadmap (ordered most to least important):
- containerize
- array of repos as json config file 
- helm chart
- tekton integration

## Configuration

| Environment Variable | Description |
| ----------- | ----------- |
| slack_token | Oauth token for Your Workspace |
| slack_channel | Channel ID to receive notifications |

## Build

```bash
go build -o ./build/releasebot
```

## Run

```bash
./build/releasebot
```
