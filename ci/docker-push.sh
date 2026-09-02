#!/usr/bin/env bash
#
# Publish one of the two images this repository owns, named by a release tag.
# Browser images need per-vendor version resolution and are built elsewhere.
#
# The tag says which image and which version, so nothing is rebuilt unless a
# release for it was cut:
#
#   ci/docker-push.sh browser-base-1.0.1
#   ci/docker-push.sh video-recorder-1.2.0
#
# A version with a pre-release suffix (1.2.0-rc1) publishes only that exact
# tag, so an rc never becomes the recorder every hub pulls by default.
#
#   DRY_RUN=1 ci/docker-push.sh browser-base-1.0.1    # print, build nothing
set -euo pipefail

TAG="${1:?usage: docker-push.sh <browser-base-VERSION|video-recorder-VERSION>}"
UBUNTU_VERSION="${UBUNTU_VERSION:-26.04}"

case "$TAG" in
    browser-base-*)   IMAGE=browser-base;   VERSION="${TAG#browser-base-}" ;;
    video-recorder-*) IMAGE=video-recorder; VERSION="${TAG#video-recorder-}" ;;
    *) echo "error: '$TAG' names no image — use browser-base-<version> or video-recorder-<version>" >&2
       exit 1 ;;
esac
[ -n "$VERSION" ] || { echo "error: '$TAG' carries no version" >&2; exit 1; }

# Docker rejects upper-case repository names; the org has capitals.
ORG=$(echo "${GITHUB_REPOSITORY:-WebSummoner}" | cut -d/ -f1 | tr '[:upper:]' '[:lower:]')
REPO="$ORG/$IMAGE"

TAGS=("$VERSION")
if [[ "$VERSION" =~ ^[0-9]+(\.[0-9]+)*$ ]]; then
    # The hub pulls websummoner/video-recorder:latest-release by default.
    [ "$IMAGE" = "video-recorder" ] && TAGS+=("latest-release")
    TAGS+=("latest")
else
    echo "==> $VERSION is a pre-release: publishing that tag only"
fi

BUILD_ARGS=()
CONTEXT=selenium/video
if [ "$IMAGE" = "browser-base" ]; then
    CONTEXT=selenium/base
    BUILD_ARGS=(--build-arg UBUNTU_VERSION="$UBUNTU_VERSION")
fi

TAG_ARGS=()
for t in "${TAGS[@]}"; do TAG_ARGS+=(-t "$REPO:$t"); done

echo "==> $REPO from $CONTEXT as: ${TAGS[*]}"

if [ "${DRY_RUN:-0}" = "1" ]; then
    echo "    docker build --pull ${BUILD_ARGS[*]} ${TAG_ARGS[*]} $CONTEXT"
    for t in "${TAGS[@]}"; do echo "    docker push $REPO:$t"; done
    exit 0
fi

[ -n "${DOCKER_USERNAME:-}" ] && [ -n "${DOCKER_PASSWORD:-}" ] || {
    echo "DOCKER_USERNAME and DOCKER_PASSWORD must be set" >&2; exit 1; }
printf '%s' "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

docker build --pull "${BUILD_ARGS[@]}" "${TAG_ARGS[@]}" "$CONTEXT"
for t in "${TAGS[@]}"; do docker push "$REPO:$t"; done
echo "==> published $REPO: ${TAGS[*]}"
