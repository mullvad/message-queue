# Golang build step / Debian Buster 21.01 / Golang 1.15.7
FROM golang:1.18.8@sha256:ceb7ae25be83b903ef598401bbbef64742fce5177df32a81bfefd4a7f0be6178 AS gobuilder
ARG version
ARG branch
ARG revision
COPY . /message-queue
WORKDIR /message-queue
RUN go install -v -ldflags="-X 'main.Branch=${branch}' -X 'main.Revision=${revision}' -X 'main.Version=${version}'" ./...

# Copy message-queue binary
FROM debian:bullseye-slim@sha256:e8ad0bc7d0ee6afd46e904780942033ab83b42b446b58efa88d31ecf3adf4678
WORKDIR /app
COPY --from=gobuilder /go/bin/message-queue .
EXPOSE 8080
CMD ["./message-queue"]
