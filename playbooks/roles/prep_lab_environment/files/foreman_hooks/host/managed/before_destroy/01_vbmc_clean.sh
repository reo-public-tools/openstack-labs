#!/bin/bash

set -x

# update, create, before_destroy etc.
HOOK_EVENT=$1

# to_s representation of the object, e.g. host's fqdn
HOOK_OBJECT=$2

# Create a temp file to dump stdin into. Remove after script is finished.
HOOK_OBJECT_FILE=$(mktemp -t foreman_hooks.XXXXXXXXXX)
echo "${HOOK_OBJECT_FILE}"
trap "rm -f $HOOK_OBJECT_FILE" EXIT
cat > $HOOK_OBJECT_FILE

# Utilize the envent data and hammer commands to pull the connectivity and bmc info.
export COMPUTE_RESOURCE=$(jq .host.compute_resource_name $HOOK_OBJECT_FILE | sed -e 's/"//g')
export VBMC_PORT=$(jq '.host.all_parameters[] | select(.name=="vbmc_port") | .value' $HOOK_OBJECT_FILE | sed -e 's/"//g')
if [ "$VBMC_PORT" == "" ]; then
  echo "Exiting as host does not have a vbmc port defined"
  exit 0
fi
COMPUTE_IP=$(/usr/bin/sudo /usr/bin/hammer compute-resource info --name ${COMPUTE_RESOURCE} --fields "Url" | awk -F '/' '{print $3}')


# Add the virtualbmc entry
if [ "${COMPUTE_IP}" != "" ]; then
  # The delete calls a stop on the instance(sigterm on the child vbmcd process). 
  # It can take 3+ minutes for the process to die after the sigterm is sent.
  # The vbmcd port for that instance stays open until the process dies.
  ssh ${COMPUTE_IP} "vbmc delete ${HOOK_OBJECT}"
fi

