#!/usr/bin/env bash

# Compile a static binary for Linux
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gbs .

# Build the container image
docker build -t gbstandings .

# Clean up
rm gbs
