# drone-gtalk

[![Build Status](https://travis-ci.org/appleboy/drone-gtalk.svg?branch=master)](https://travis-ci.org/appleboy/drone-gtalk) [![codecov](https://codecov.io/gh/appleboy/drone-gtalk/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-gtalk) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-gtalk)](https://goreportcard.com/report/github.com/appleboy/drone-gtalk)

Drone plugin for sending Gtalk notifications.

## Build

Build the binary with the following commands:

```
$ make build
```

## Testing

Test the package with the following command:

```
$ make test
```

## Docker

Build the docker image with the following commands:

```
$ make docker
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-gtalk' not found or does not exist..
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e PLUGIN_FB_PAGE_TOKEN=xxxxxxx \
  -e PLUGIN_FB_VERIFY_TOKEN=xxxxxxx \
  -e PLUGIN_TO=xxxxxxx \
  -e PLUGIN_MESSAGE=test \
  -e DRONE_REPO_OWNER=appleboy \
  -e DRONE_REPO_NAME=go-hello \
  -e DRONE_COMMIT_SHA=e5e82b5eb3737205c25955dcc3dcacc839b7be52 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=appleboy \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/appleboy/go-hello \
  appleboy/drone-gtalk
```
