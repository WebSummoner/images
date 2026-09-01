#!/bin/bash

set -e

export GO111MODULE="on"
LDFLAGS="-X github.com/websummoner/images/cmd.buildStamp=$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X github.com/websummoner/images/cmd.gitRevision=$(git describe --tags || git rev-parse HEAD) -s -w"

build() {
    local goos=$1 goarch=$2 ext=""
    if [ "$goos" = "windows" ]; then
        ext=".exe"
    fi
    GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o "dist/images_${goos}_${goarch}${ext}" .
}

build linux amd64
build darwin amd64
build windows amd64
build windows 386
