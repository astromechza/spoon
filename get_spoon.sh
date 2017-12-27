#!/usr/bin/env sh

set -eu
set -o pipefail

# destination dir can be overriden
DESTINATION_DIR=${DESTINATION_DIR:-/usr/bin/}

PLATFORM=""
UNAME=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "${UNAME}" in
    linux*) PLATFORM="linux" ;;
    darwin*) PLATFORM="darwin" ;;
    *) echo "Unsupported platform ${UNAME}"; exit 1;;
esac
case "${ARCH}" in
    x86_64) PLATFORM="${PLATFORM}-amd64" ;;
    *) echo "Unsupported architecture ${ARCH}" ;;
esac

# construct the final artifact name
ARTIFACT="spoon-${PLATFORM}"

# pull latest release json
echo "Finding latest release.."
LATEST_RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/AstromechZA/spoon/releases/latest)

# isolate tag name
RELEASE_TAG=$(echo "${LATEST_RELEASE}" | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
echo "Got ${RELEASE_TAG}"

# downloading artifact
echo "Downloading ${ARTIFACT} to ${DESTINATION_DIR}.."
curl -L \
    -o "${DESTINATION_DIR}/spoon" \
    "https://github.com/AstromechZA/spoon/releases/download/${RELEASE_TAG}/${ARTIFACT}"
