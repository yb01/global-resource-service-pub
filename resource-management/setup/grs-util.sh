#!/usr/bin/env bash

GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

CLOUD_PROVIDER="${CLOUD_PROVIDER:-gce}"

# PROVIDER_VARS is a list of cloud provider specific variables. Note:
# this is a list of the _names_ of the variables, not the value of the
# variables.

PROVIDER_UTILS="${GRS_ROOT}/setup/${CLOUD_PROVIDER}/util.sh"
if [ -f "${PROVIDER_UTILS}" ]; then
    source "${PROVIDER_UTILS}"
fi
