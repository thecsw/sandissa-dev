#!/bin/sh

GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

export GOOS
export GOARCH
export CGO_ENABLED

go build -v
