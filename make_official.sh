#!/usr/bin/env bash

set -e

mkdir -p build/

# add the git commit id and date
VERSION="$(cat VERSION) (commit $(git rev-parse --short HEAD) @ $(git log -1 --date=short --pretty=format:%cd))"

function buildbinary {
    export GOOS=$1
    export GOARCH=$2

    echo "Building official ${GOOS} ${GOARCH} binary for version '${VERSION}'"

    go build -i -v -o "build/spoon-${GOOS}-${GOARCH}" -ldflags "-X \"main.SpoonVersion=${VERSION}\""

    echo "Done"
    ls -l "build/spoon-${GOOS}-${GOARCH}"
    file "build/spoon-${GOOS}-${GOARCH}"
    echo

    unset GOOS
    unset GOARCH
}

# platform builds
buildbinary darwin amd64
buildbinary linux amd64

# and build for dev
go build -ldflags "-X \"main.SpoonVersion=${VERSION}\""
