#!/bin/bash

# Array of (destination,interface) tuples
# If new element needs to be added, please delimit destination and interface by comma like other elements
array=( anonymous,Anonymous
        command,Handler 
        plugin,Plugin
        store,Store
        utils/store,KVStore)
        
for i in "${array[@]}"; do IFS=",";
    set -- ${i};

    DESTINATION=${1}
    INTERFACE=${2}
    
    # Test echoing
    echo "${DESTINATION} ${INTERFACE}"
done