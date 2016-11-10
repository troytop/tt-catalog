#!/usr/bin/env bash -x

: ${vars_file="../configs/vars.yml"}
: ${target_alias="catalyst"}
: ${target="https://52.160.107.255:8081"}

CURRENT_DIR=`pwd`
SCRIPT_DIR_NAME=`dirname "$0"`

CI_DIR=`echo "${CURRENT_DIR}/${SCRIPT_DIR_NAME}"`

export gpg_private_key=$(cat ${GPG_KEY_DEVELOPMENT})
export gpg_key_passphrase=${PASSPHRASE_DEVELOPMENT}
export pipeline_file="${CI_DIR}/catalog-publisher.yml"
export pipeline_name="catalog-publisher-development"

fly -t $target_alias set-pipeline -p $pipeline_name -c ${pipeline_file} --load-vars-from ${vars_file} --var "git-private-key=$(cat ~/.ssh/hsm.private.shared.key)" --var "gpg_private_key=${gpg_private_key}" --var "gpg_key_passphrase=${gpg_key_passphrase}"

fly -t $target_alias unpause-pipeline -p $pipeline_name

open ${target}/pipelines/$pipeline_name
