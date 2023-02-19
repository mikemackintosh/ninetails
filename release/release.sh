#!/bin/bash

PACKAGE="github.com/mikemackintosh/ninetails/internal"

# ARCHITECTURES SUPPORTED
ARCH=(
    amd64
    arm64
)

# OPERATING SYSTEMS SUPPORT
OS=(
    linux
    darwin
    windows
)

git fetch --tags
VERSION="$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

# STEP 2: Build the ldflags

LDFLAGS=(
  "-X '${PACKAGE}/version.Version=${VERSION}'"
  "-X '${PACKAGE}/version.CommitHash=${COMMIT_HASH}'"
  "-X '${PACKAGE}/version.BuildTime=${BUILD_TIMESTAMP}'"
)

# STEP 3: Actual Go build process

go build -ldflags="${LDFLAGS[*]}"

if [[ ! -d release/bin ]]; then 
    mkdir release/bin
fi

for os in ${OS[@]}; do 
    for arch in ${ARCH[@]}; do
        echo -e "\nBuilding for $os-$arch"
        GOOS=$os GOARCH=$arch go build -ldflags="${LDFLAGS[*]}" -o release/bin/ninetails-$os-$arch cmd/main.go
        echo $(shasum -a 256 release/bin/ninetails-$os-$arch)
    done
done