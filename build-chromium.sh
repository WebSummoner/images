#!/bin/bash

VERSION=${1:-101.0.4951.64-0ubuntu0.18.04.1}
TAG=${2:-chromium_101.0}
BASE_TAG=${3:-1.0.0}

# Cleanup stuff
export BUILDKIT_PROGRESS=plain
docker rmi -f websummoner/vnc:$TAG websummoner/browser-base:$BASE_TAG $(docker images -q websummoner/dev_chromium:*)
rm -rf ../websummoner-container-tests

# Prepare for building images
go build

# Forked tests with a bugfix
git clone -b add-missing-dependency https://github.com/sskorol/websummoner-container-tests.git ../websummoner-container-tests

# Force build websummoner/browser-base image as it has arm64-specific updates
cd ./selenium/base && docker build --no-cache --build-arg UBUNTU_VERSION=18.04 -t websummoner/browser-base:$BASE_TAG . && docker system prune -f

# Build chromium image
cd ../../ && ./images chromium -b $VERSION -t websummoner/vnc:$TAG --test && docker system prune -f
