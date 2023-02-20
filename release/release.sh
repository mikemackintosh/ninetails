#!/bin/bash

if [[ -z $NINETAILS_GITHUB_TOKEN ]]; then 
    echo "Please set 'NINETAILS_GITHUB_TOKEN'"
    exit 1
fi

# Load common vars
source ./release/common.sh
if [[ -z $VERSION || -z $COMMIT_HASH ]]; then
    echo "VERSION or COMMIT_HASH unset. Exiting."
    exit 1
fi

# Create next tag
CREATE_TAG_OUTPUT=$(curl \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${NINETAILS_GITHUB_TOKEN}"\
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/${REPO}/git/tags \
  -d '{"tag":"'$VERSION'","message":"'$VERSION'","object":"'$COMMIT_HASH'","type":"commit"}')
echo $CREATE_TAG_OUTPUT

exit 0
# Build
./release/build

# Add binaries to next tag

# Create release