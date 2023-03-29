# Golang build step / Debian Buster 21.01 / Golang 1.15.7
FROM docker.io/library/golang:1.19.7@sha256:d2078d2176fcc81c7df4dbd5c5859147474a5e46bd9263460f05b820019e5fcc AS gobuilder
ARG version
ARG branch
ARG revision
COPY . /message-queue
WORKDIR /message-queue
RUN go install -v -ldflags="-X 'main.Branch=${branch}' -X 'main.Revision=${revision}' -X 'main.Version=${version}'" ./...

# Copy message-queue binary
FROM docker.io/library/debian:bullseye-slim@sha256:7acda01e55b086181a6fa596941503648e423091ca563258e2c1657d140355b1
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
