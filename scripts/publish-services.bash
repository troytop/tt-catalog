#!/bin/bash -x

#####
# Script used by CSM nightly pipeline to publish the service SDLs to build dev-catalog bucket
# Bucket location is specified by S3_CATALOG_LOCATION variable
# REGISTRY_LOCATION specifies the docker registry that is being used
#####

OS_VERSION=$(go env GOOS)

: ${DOCKER_REPOSITORY:=docker.io}
: ${DEFAULT_IMAGE_TAG:=latest}
: ${DEFAULT_SDL_VERSION:=1.0.0}
: ${SERVICE_DIR:=services/HPE}
: ${S3_CATALOG_LOCATION:=s3://helion-service-manager/build/dev-catalog}

if [ -z ${APP_VERSION_TAG} ]; then
	echo "Error: APP_VERSION_TAG is not set"
	exit 1
fi


if [ -z ${SDL_VERSION} ]; then
	echo "Error: SDL_VERSION is not set"
	exit 1
fi

# uploads files to s3 catalog
function upload_file_to_s3_catalog(){
	# $1 - source (local file path)
	# $2 - destination (s3 location)
	aws s3 cp $1 $2 --acl public-read
}

# publish the SDL file to the s3 catalog
function publish_service_sdl(){
	for sdl_version in `ls -l $1 | grep "^d" | awk -F" " '{print $9}'`
	do
		file_path=$1/${sdl_version}/sdl/sdl.json
		echo $file_path
		#darwin
		if [ "$OS_VERSION" = "darwin" ]; then
			sed -i '.original' "\|repository|s|docker.io|${DOCKER_REPOSITORY}|gwp" ${file_path}
			sed -i '.updated_repository' "\|tag|s|${DEFAULT_IMAGE_TAG}|${APP_VERSION_TAG}|gwp" ${file_path}
			sed -i '.updated_sdl_version' "\|sdl_version|s|${DEFAULT_SDL_VERSION}|${SDL_VERSION}|gwp" ${file_path}
		elif [ "$OS_VERSION" = "linux" ]; then
			sed -i "\|repository|s|docker.io|${DOCKER_REPOSITORY}|gwp" ${file_path}
			sed -i "\|tag|s|${DEFAULT_IMAGE_TAG}|${APP_VERSION_TAG}|g" ${file_path}
			sed -i "\|sdl_version|s|${DEFAULT_SDL_VERSION}|${SDL_VERSION}|gwp" ${file_path}
		fi

		# check if file is updated
		grep tag $file_path | grep ${APP_VERSION_TAG}
		if [ $? -ne 0 ]; then
			echo "Error: Image tag is not updated"
			exit 1
		fi

		# check if file is updated
		grep sdl_version $file_path | grep ${SDL_VERSION}
		if [ $? -ne 0 ]; then
			echo "Error: sdl_version is not updated"
			exit 1
		fi

		upload_file_to_s3_catalog $file_path ${S3_CATALOG_LOCATION}/$1/${SDL_VERSION}/sdl/
		
		# upload sdl_version config 
		if [ -f $1/${sdl_version}/config ]
		then
			upload_file_to_s3_catalog $1/${sdl_version}/config ${S3_CATALOG_LOCATION}/$1/${SDL_VERSION}/
		fi
	done
}

# 
function publish_service(){
	for product_version in `ls -l $1 | grep "^d" | awk -F" " '{print $9}'`
	do
		publish_service_sdl $1/$product_version
	done

	for catalog_file in `ls -l $1 | grep "^-" | awk -F" " '{print $9}'`
	do
		upload_file_to_s3_catalog $1/$catalog_file ${S3_CATALOG_LOCATION}/$1/$catalog_file
	done

}

for id in `ls ${SERVICE_DIR}`
do
	if [[ $SKIP_SDL_VALIDATION =~ $id ]]; then
		echo "SKIPPED"
	else
		publish_service "${SERVICE_DIR}/$id"
	fi
done