# message-queue
A message queue that reads messages from redis pubsub, and publishes them to clients connected via websocket
## Building

Clone this repository, and run `make` to build.
This will produce a `message-queue` binary and put them in your `GOBIN`.

## Testing
To run the tests, run `make test`.
To run the integration tests as well, run `go test ./...`. Note that this requires a local instance of redis and redis-sentinel.

## Run In Browser
[![Open with Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/from-referrer/)

## Usage
All options can be either configured via command line flags, or via their respective environment variable, as denoted by `[ENVIRONMENT_VARIABLE]`.
To get a list of all the options, run `message-queue -h`.

## Packaging
In order to deploy message-queue, we use docker.

To build a new image:
- Update the version in VERSION
- run `make package`.

This will create a new image tagged as `quay.io/mullvad/message-queue:<version>`.

Current docker repo digests:

|   tag    |                                             repo path                                             |
|:--------:|:-------------------------------------------------------------------------------------------------:|
| `1.0.0-buster-21.01`  | `quay.io/mullvad/message-queue@sha256:e8da7429612b7954732d1bc19f3d828a7ca193f676398a4d2432130c35eb1406` |
| `1.0.0-buster-20.09`  | `quay.io/mullvad/message-queue@sha256:d319005c398ee068afc0967030e56a9c2d4515d52e65440c26a8e17c89e216ba` |
| `1.0.0`  | `quay.io/mullvad/message-queue@sha256:8722d0ba33eec49c809af017f70c56cbfe6593b4b0b05c3e13012776f610b590` |
