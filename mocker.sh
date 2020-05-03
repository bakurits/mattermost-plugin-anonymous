#!/bin/bash

# Github Repository Link
GH_REPO=github.com/bakurits/mattermost-plugin-anonymous

# Array of (package,interface) tuples
# If new element needs to be added, please delimit package and interface by colon like other elements
array=(
    anonymous:Anonymous
    command:Handler
    plugin:Plugin
    store:Store
    utils/store:API
)

# Iterate over each item in array and generate mock
for i in "${array[@]}"; do IFS=":";
    set -- ${i};

    PACKAGE=server/${1}
    INTERFACE=${2}
    MOCK_FILE=$(echo ${INTERFACE} | awk '{print tolower($0)}')_mock.go

    echo "Generating Mock for Interaface ${INTERFACE} in /${PACKAGE}..."

    # Actual Command
    mockgen -destination=${PACKAGE}/mock/${MOCK_FILE} -package=mock ${GH_REPO}/${PACKAGE} ${INTERFACE}
done
