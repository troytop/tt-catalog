#!/bin/bash

if [ -z ${CATALOG_LOCATION_DEVELOPMENT} ]; then
    echo "Error: CATALOG_LOCATION_DEVELOPMENT is not set. It is recommended to run this script from `make publish-catalog`"
    exit 1
fi

if [ -z ${IDL_LOCATION} ]; then
    echo "Error: CATALOG_LOCATION is not set. It is recommended to run this script from `make publish-catalog`"
    exit 1
fi

if [ -d services ]; 
then
    aws s3 sync services ${CATALOG_LOCATION_DEVELOPMENT} --exclude "*/*.sig" --acl public-read --delete
fi

if [ -d instance-definition ];
then
    aws s3 sync instance-definition ${IDL_LOCATION} --acl public-read --delete
fi

BASEPATH=temp
rm -rf $BASEPATH
mkdir -p $BASEPATH
aws s3 cp ${CATALOG_LOCATION_DEVELOPMENT} temp --recursive --exclude "*" --include "*/sdl.*" > /dev/null
for i in $(find temp -name sdl.sig);
do
  DIR=$(dirname $i)
  if [ -z $(find $DIR -name sdl.json) ]; then
    FILEPATH=$(echo $DIR | sed 's|'$BASEPATH'/|/|')
    SDL_DIR=$(dirname $FILEPATH)
    aws s3 rm ${CATALOG_LOCATION_DEVELOPMENT}${SDL_DIR} --recursive
  fi
done
