#!/bin/bash

# Github Repository Link
GH_REPO=github.com/bakurits/mattermost-plugin-anonymous

# Array of (destination,interface) tuples
# If new element needs to be added, please delimit destination and interface by comma like other elements
array=(
    anonymous,Anonymous
    command,Handler
    plugin,Plugin
    store,Store
    utils/store,KVStore
)

# Iterate over each item in array and generate mock
for i in "${array[@]}"; do IFS=",";
    set -- ${i};

    DESTINATION=server/${1}
    INTERFACE=${2}
    MOCK_FILE=$(echo ${INTERFACE} | awk '{print tolower($0)}')_mock.go

    echo "Generating Mock for Interaface ${INTERFACE} in /${DESTINATION}..."

    # Actual Command
    mockgen -destination=${DESTINATION}/mock/${MOCK_FILE} -package=mock ${GH_REPO}/${DESTINATION} ${INTERFACE}
done