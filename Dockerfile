# Copyright 2018 Platform 9 Systems, Inc.

# Reproducible builder image
FROM golang:1.10.0 as builder
WORKDIR /go/src/github.com/platform9/ssh-provider
# This expects that the context passed to the docker build command is
# the ssh-provider directory.
# e.g. docker build -t <tag> -f <this_Dockerfile> <path_to_cluster-api>
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install -a -ldflags '-extldflags "-static"' github.com/platform9/ssh-provider

# Final container
FROM debian:stretch-slim
RUN apt-get update && apt-get install -y ca-certificates openssh-server && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/ssh-provider .
