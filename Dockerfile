# Golang build step / Debian Buster 21.01 / Golang 1.15.7
FROM docker.io/library/golang:1.18.10@sha256:50c889275d26f816b5314fc99f55425fa76b18fcaf16af255f5d57f09e1f48da AS gobuilder
ARG version
ARG branch
ARG revision
COPY . /message-queue
WORKDIR /message-queue
RUN go install -v -ldflags="-X 'main.Branch=${branch}' -X 'main.Revision=${revision}' -X 'main.Version=${version}'" ./...

# Copy message-queue binary
FROM docker.io/library/debian:bullseye-slim@sha256:d51d5c391d202d5e2e0294a9df6ff077ed40583b11831d347d418690da496c50
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=gobuilder /go/bin/message-queue .
EXPOSE 8080
CMD ["./message-queue"]
