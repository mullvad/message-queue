# Golang build step / Debian Buster 21.01 / Golang 1.15.7
FROM golang:1.17.6@sha256:8fb6a244eeea71ca938a47eb046a39999a96ad43a0a69bf4275463df135bdc95 AS gobuilder
ARG version
ARG branch
ARG revision
COPY . /message-queue
WORKDIR /message-queue
RUN go install -v -ldflags="-X 'main.Branch=${branch}' -X 'main.Revision=${revision}' -X 'main.Version=${version}'" ./...

# Copy message-queue binary
FROM debian:stretch@sha256:c5c5200ff1e9c73ffbf188b4a67eb1c91531b644856b4aefe86a58d2f0cb05be
WORKDIR /app
COPY --from=gobuilder /go/bin/message-queue .
EXPOSE 8080
CMD ["./message-queue"]
