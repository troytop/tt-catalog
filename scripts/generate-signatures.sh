#!/bin/sh

if [ -z ${BUCKET} ]; then
    echo "Error: BUCKET is not set. It is recommended to run this script from 'make sign-development' or 'make sign-stable'"
    exit 1
fi

if [ -z ${PRIVATE_KEY} ]; then
    echo "Error: PRIVATE_KEY is not set. It is recommended to run this script from 'make sign-development' or 'make sign-stable"
    exit 1
fi

if [ -z ${PASSPHRASE} ]; then
    echo "Error: PASSPHRASE is not set. It is recommended to run this script from 'make sign-development' or 'make sign-stable"
    exit 1
fi

if [ -z $(pip freeze | grep awscli) ]; then
  echo "awscli is not available please install with 'pip install awscli'"
  exit 1
fi
BASEPATH=temp

mkdir $BASEPATH
aws s3 cp ${BUCKET} temp --recursive --exclude "*" --include "*/sdl.*" > null


for i in $(find temp -name sdl.json);
do
  DIR=$(dirname $i)
  
  if [ -z $(find $DIR -name sdl.sig) ]; then
    FILEPATH=$(echo $DIR | sed 's|'$BASEPATH'/|/|')
    echo $FILEPATH
    echo $PASSPHRASE | gpg --detach-sign --passphrase-fd 0 -o $DIR/sdl.sig --default-key $PRIVATE_KEY $i
    
    aws s3 cp ${DIR}/sdl.sig ${BUCKET}${FILEPATH}/sdl.sig --acl public-read
  fi
done

rm -r temp
