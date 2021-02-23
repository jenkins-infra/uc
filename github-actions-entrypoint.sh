#!/bin/sh -le

ARGS="--display-updates --write --path ${PLUGINS_PATH}"

if [ "${SECURITY_UPDATES}" = "true" ]; then 
  ARGS="$ARGS --security-updates"
fi

if [ -n  "${JENKINS_VERSION}" ]; then
  ARGS="$ARGS --jenkins-version ${JENKINS_VERSION}"
else
  ARGS="$ARGS --determine-version-from-dockerfile"
fi

uc update $ARGS
