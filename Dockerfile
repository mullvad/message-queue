# Golang build step / Debian Buster 21.01 / Golang 1.15.7
FROM docker.io/library/golang:1.18.9@sha256:4833a2f09eab2c1dada9b5e43078e6c24572b62f74db913d1839ae8c6fe7c99d AS gobuilder
ARG version
ARG branch
ARG revision
COPY . /message-queue
WORKDIR /message-queue
RUN go install -v -ldflags="-X 'main.Branch=${branch}' -X 'main.Revision=${revision}' -X 'main.Version=${version}'" ./...

# Copy message-queue binary
FROM docker.io/library/debian:bullseye-slim@sha256:171530d298096f0697da36b3324182e872db77c66452b85783ea893680cc1b62
WORKDIR /app
COPY --from=gobuilder /go/bin/message-queue .
EXPOSE 8080
CMD ["./message-queue"]
