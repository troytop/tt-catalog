#!/bin/bash

if [ -z $1 ]; then
    # used by build pipelines
    if [ -z ${PASSPHRASE} ]; then
        echo "Error: Supply the key passphase. Use 'make sign-development PASSPHRASE=<Key Passphrase>'. Or set PASSPHRASE environment variable"
        exit 1
    fi
else
    export PASSPHRASE=$1
fi

export BUCKET=${CATALOG_LOCATION_DEVELOPMENT}
export PRIVATE_KEY=87E97E85

scripts/generate-signatures.sh