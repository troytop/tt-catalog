#!/bin/bash

if [ -z $1 ]; then
    echo "Error: Supply the key passphase. Use 'make sign-stable PASSPHRASE=<Key Passphrase>'"
    exit 1
fi

export BUCKET=${CATALOG_LOCATION_STABLE}
export PRIVATE_KEY=573A2CCA
export PASSPHRASE=$1

scripts/generate-signatures.sh