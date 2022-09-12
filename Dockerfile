# Golang build step / Debian Buster 21.01 / Golang 1.15.7
FROM golang:1.18.6@sha256:f1741ed3855b1a558f3c057f777667a4b73ff250c4f1bee7f0388cc2d1224079 AS gobuilder
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
