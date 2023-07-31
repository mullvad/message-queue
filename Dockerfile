# Golang build step / Debian Buster 21.01 / Golang 1.15.7
FROM docker.io/library/golang:1.19.11@sha256:e8579154cca7ee15a5ea89ba7ad3cc3967253c2cd3010fd28fa0573b8d06c839 AS gobuilder
ARG version
ARG branch
ARG revision
COPY . /message-queue
WORKDIR /message-queue
RUN go install -v -ldflags="-X 'main.Branch=${branch}' -X 'main.Revision=${revision}' -X 'main.Version=${version}'" ./...

# Copy message-queue binary
FROM docker.io/library/debian:bullseye-slim@sha256:fd3b382990294beb46aa7549edb9f40b11a070f959365ef7f316724b2e425f90
# Same ARGs as in the first stage to set labels in the final image
ARG version
ARG branch
ARG revision
LABEL org.opencontainers.image.version="$version" org.opencontainers.image.ref.name="$branch" org.opencontainers.image.revision="$revision"
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=gobuilder /go/bin/message-queue .
EXPOSE 8080
CMD ["./message-queue"]
