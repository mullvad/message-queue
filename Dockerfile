# Golang build step / Debian Buster 21.01 / Golang 1.15.7
FROM golang:1.19.0@sha256:d72a2944621f0f4027983e094bdf9464d26459c17efe732320043b84c91e51e7 AS gobuilder
ARG version
ARG branch
ARG revision
COPY . /message-queue
WORKDIR /message-queue
RUN go install -v -ldflags="-X 'main.Branch=${branch}' -X 'main.Revision=${revision}' -X 'main.Version=${version}'" ./...

# Copy message-queue binary
FROM debian:stretch@sha256:4bb600434787c903886fe33526d19ff33114a33b433a4a4cdbdf9b8543f1ab5d
WORKDIR /app
COPY --from=gobuilder /go/bin/message-queue .
EXPOSE 8080
CMD ["./message-queue"]
