#!/usr/out/env bash
set -eu

if [ "${PARAM_VERBOSE}" = "true" ] ; then
  set -x
fi

EXT_LABEL_1="io.buildpack.extension.layers"
EXT_LABEL_2="io.buildpack.templates.order-extensions"

IMG_MANIFEST=$(skopeo inspect --authfile ${PARAM_USER_HOME}/creds-secrets/dockercfg/.dockerconfigjson "docker://${PARAM_BUILDER_IMAGE}")