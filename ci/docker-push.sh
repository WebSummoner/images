#!/usr/bin/env bash
#
# Publish the two images this repository owns: browser-base and the video
# recorder. Browser images need per-vendor version resolution and are built
# elsewhere.
#
#   ci/docker-push.sh latest
#   BASE_TAG=1.0.1 ci/docker-push.sh 1.0.1
set -e

MOVING_TAG="${1:?usage: docker-push.sh <tag>}"
BASE_TAG="${BASE_TAG:-1.0.0}"
UBUNTU_VERSION="${UBUNTU_VERSION:-26.04}"

# Docker rejects upper-case repository names; the org has capitals.
ORG=$(echo "${GITHUB_REPOSITORY%%/*}" | tr '[:upper:]' '[:lower:]')

[ -n "${DOCKER_USERNAME:-}" ] && [ -n "${DOCKER_PASSWORD:-}" ] || {
    echo "DOCKER_USERNAME and DOCKER_PASSWORD must be set" >&2; exit 1; }
printf '%s' "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

echo "==> $ORG/browser-base:$BASE_TAG"
docker build --pull --build-arg UBUNTU_VERSION="$UBUNTU_VERSION" \
    -t "$ORG/browser-base:$BASE_TAG" -t "$ORG/browser-base:$MOVING_TAG" selenium/base
docker push "$ORG/browser-base:$BASE_TAG"
docker push "$ORG/browser-base:$MOVING_TAG"

# The hub defaults to :latest-release.
echo "==> $ORG/video-recorder"
docker build --pull \
    -t "$ORG/video-recorder:latest-release" -t "$ORG/video-recorder:$MOVING_TAG" selenium/video
docker push "$ORG/video-recorder:latest-release"
docker push "$ORG/video-recorder:$MOVING_TAG"

echo "==> published"
