#!/usr/bin/env bash

# Bring timezone info into the local directory
cp /usr/local/go/lib/time/zoneinfo.zip .

# Compile a static binary for Linux
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gbs .

# Build the container image
docker build -t gbstats .

# Clean up
rm gbs
rm zoneinfo.zip

