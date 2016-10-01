#!/bin/bash

if [ -z $1 ]; then
    echo "Error: Supply the key passphase. Use 'make sign-development PASSPHRASE=<Key Passphrase>'"
    exit 1
fi

export BUCKET=${CATALOG_LOCATION_DEVELOPMENT}
export PRIVATE_KEY=87E97E85
export PASSPHRASE=$1

scripts/generate-signatures.sh