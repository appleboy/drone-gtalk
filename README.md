# drone-gtalk

![logo](images/logo.png)

[![GoDoc](https://godoc.org/github.com/appleboy/drone-gtalk?status.svg)](https://godoc.org/github.com/appleboy/drone-gtalk)
[![codecov](https://codecov.io/gh/appleboy/drone-gtalk/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-gtalk)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-gtalk)](https://goreportcard.com/report/github.com/appleboy/drone-gtalk)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-gtalk.svg)](https://hub.docker.com/r/appleboy/drone-gtalk/)
[![microbadger](https://images.microbadger.com/badges/image/appleboy/drone-gtalk.svg)](https://microbadger.com/images/appleboy/drone-gtalk "Get your own image badge on microbadger.com")
[![Build status](https://ci.appveyor.com/api/projects/status/d7t9jb5ouoa7tk6i?svg=true)](https://ci.appveyor.com/project/appleboy/drone-gtalk)

[Drone](https://github.com/drone/drone) plugin for sending Gtalk notifications.

## Build or Download a binary

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/drone-gtalk/releases). Support the following OS type.

* Windows amd64/386
* Linux amd64/386
* Darwin amd64/386

With `Go` installed

```sh
go get -u -v github.com/appleboy/drone-gtalk
```

or build the binary with the following command:

```sh
make build
```

## Testing

Test the package with the following command:

```sh
make test
```

## Docker

Build the docker image with the following commands:

```sh
make docker
```

## Usage

Execute from the working directory:

```bash
docker run --rm \
  -e PLUGIN_HOST=talk.google.com:443 \
  -e PLUGIN_USERNAME=xxxxxxx \
  -e PLUGIN_OAUTH_TOKEN=xxxxxxx \
  -e PLUGIN_TO=xxxxxxx \
  -e PLUGIN_MESSAGE=test \
  -e PLUGIN_ONLY_MATCH_EMAIL=false \
  -e DRONE_REPO_OWNER=appleboy \
  -e DRONE_REPO_NAME=go-hello \
  -e DRONE_COMMIT_SHA=e5e82b5eb3737205c25955dcc3dcacc839b7be52 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=appleboy \
  -e DRONE_COMMIT_AUTHOR_EMAIL=appleboy@gmail.com \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/appleboy/go-hello \
  -e DRONE_JOB_STARTED=1477550550 \
  -e DRONE_JOB_FINISHED=1477550750 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-gtalk
```

Load all environments from file.

```bash
docker run --rm \
  -e ENV_FILE=your_env_file_path \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-gtalk
```
