#!/usr/bin/env sh

set -eu
set -o pipefail

# destination dir can be overriden
DESTINATION_DIR=${DESTINATION_DIR:-/usr/bin/}

ARTIFACT=binaries.tar.gz

# pull latest release json
echo "Finding latest release.."
LATEST_RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/AstromechZA/spoon/releases/latest)

# isolate tag name
RELEASE_TAG=$(echo "${LATEST_RELEASE}" | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
echo "Got ${RELEASE_TAG}."

# make temp directory
tempdir="$(mktemp -d 2>/dev/null || mktemp -d -t 'spoon-tmp')"
echo "Created tempdir ${tempdir}."

# download it
echo "Downloading ${ARTIFACT}.."
curl -L \
    -o "${tempdir}/${ARTIFACT}" \
    "https://github.com/AstromechZA/spoon/releases/download/${RELEASE_TAG}/${ARTIFACT}"

# untar it
tar -xzvf "${tempdir}/${ARTIFACT}" -C "${tempdir}"

# identify platform dir
platform=$(uname -a)
if [ "$(echo "${platform}" | grep -ic 'linux')" -ge 1 ]; then
    platform="linux_amd64"
elif [ "$(echo "${platform}" | grep -ic 'darwin')" -ge 1 ]; then
    platform="darwin_amd64"
else
    echo "Error: Unable to detect compatible platform (${platform})."
    exit 1
fi

echo "Copying ${tempdir}/${platform}/spoon to ${DESTINATION_DIR}.."
cp "${tempdir}/${platform}/spoon" "${DESTINATION_DIR}"

# cleanup
rm -rf "${tempdir}"
