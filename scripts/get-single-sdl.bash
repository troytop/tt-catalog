#!/bin/bash

SRC_CATALOG=s3://helion-service-manager/build/dev-catalog


if [ "$#" -ne 3 ]; then
  echo "Usage: $0 <SERVICE-NAME> <PRODUCT-VERSION> <SDL-VERSION>" >&2
  exit 1
fi

SERVICE_NAME=$1
PRODUCT_VERSION=$2
SDL_VERSION=$3

mkdir -p services/HPE/${SERVICE_NAME}/${PRODUCT_VERSION}/${SDL_VERSION}

aws s3 cp ${SRC_CATALOG}/services/HPE/${SERVICE_NAME}/config.json services/HPE/$SERVICE_NAME/config.json
aws s3 cp ${SRC_CATALOG}/services/HPE/${SERVICE_NAME}/icon.png services/HPE/$SERVICE_NAME/icon.png
aws s3 cp ${SRC_CATALOG}/services/HPE/${SERVICE_NAME}/README.md services/HPE/$SERVICE_NAME/README.md
aws s3 cp ${SRC_CATALOG}/services/HPE/${SERVICE_NAME}/${PRODUCT_VERSION}/${SDL_VERSION} services/HPE/$SERVICE_NAME/${PRODUCT_VERSION}/${SDL_VERSION} --recursive
