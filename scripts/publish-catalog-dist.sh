#!/bin/bash

export NOW=$(date +"%m%d%Y%H%M%S")

BASEPATH=$(mktemp -d)/$VERSION
rm -rf $BASEPATH
mkdir -p $BASEPATH
echo $BASEPATH
aws s3 cp ${CATALOG_LOCATION_DEVELOPMENT} ${BASEPATH}/services --recursive > /dev/null

cd ${BASEPATH}

tar -czf catalog-templates-dist.tgz services
aws s3 cp catalog-templates-dist.tgz ${CATALOG_DIST_LOCATION_DEVELOPMENT}/catalog-templates-dist.tgz --acl public-read
aws s3 cp catalog-templates-dist.tgz ${CATALOG_DIST_LOCATION_DEVELOPMENT}/${STACKATO_VERSION}/catalog-templates-dist.tgz --acl public-read
aws s3 cp catalog-templates-dist.tgz ${CATALOG_DIST_LOCATION_DEVELOPMENT}/${STACKATO_VERSION}/catalog-templates-dist-${NOW}.tgz --acl public-read

rm -rf ${BASEPATH}
