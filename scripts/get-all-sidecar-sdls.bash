#!/bin/bash -x

SRC_CATALOG=s3://helion-service-manager/build/dev-catalog


if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <SDL-VERSION>" >&2
  exit 1
fi

NEW_SDL_VERSION=$1

function get_service_definitions(){
    SERVICE_NAME=$1
    PRODUCT_VERSION=$2
    SDL_VERSION=$3

    mkdir -p services/HPE/${SERVICE_NAME}/${PRODUCT_VERSION}/${SDL_VERSION}

    aws s3 cp ${SRC_CATALOG}/services/HPE/${SERVICE_NAME}/config.json services/HPE/$SERVICE_NAME/config.json
    aws s3 cp ${SRC_CATALOG}/services/HPE/${SERVICE_NAME}/icon.png services/HPE/$SERVICE_NAME/icon.png
    aws s3 cp ${SRC_CATALOG}/services/HPE/${SERVICE_NAME}/README.md services/HPE/$SERVICE_NAME/README.md
    aws s3 cp ${SRC_CATALOG}/services/HPE/${SERVICE_NAME}/${PRODUCT_VERSION}/${SDL_VERSION} services/HPE/$SERVICE_NAME/${PRODUCT_VERSION}/${SDL_VERSION} --recursive

}

get_service_definitions mysql 5.5 ${NEW_SDL_VERSION}
get_service_definitions mongo 3.0 ${NEW_SDL_VERSION}
get_service_definitions postgres 9.4 ${NEW_SDL_VERSION}
get_service_definitions redis 3.0 ${NEW_SDL_VERSION}
get_service_definitions rabbitmq 3.6 ${NEW_SDL_VERSION}
get_service_definitions mssql 2014 ${NEW_SDL_VERSION}
get_service_definitions rds-mysql 5.5 ${NEW_SDL_VERSION}
get_service_definitions rds-postgres 9.4 ${NEW_SDL_VERSION}
get_service_definitions havenondemand 1.0 ${NEW_SDL_VERSION}
get_service_definitions pass-through 1.0 ${NEW_SDL_VERSION}
