# Golang build step / Debian Buster 21.01 / Golang 1.15.7
FROM docker.io/library/golang:1.19.8@sha256:9f2dd04486e84eec72d945b077d568976981d9afed8b4e2aeb08f7ab739292b3 AS gobuilder
ARG version
ARG branch
ARG revision
COPY . /message-queue
WORKDIR /message-queue
RUN go install -v -ldflags="-X 'main.Branch=${branch}' -X 'main.Revision=${revision}' -X 'main.Version=${version}'" ./...

# Copy message-queue binary
FROM docker.io/library/debian:bullseye-slim@sha256:f4da3f9b18fc242b739807a0fb3e77747f644f2fb3f67f4403fafce2286b431a
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
