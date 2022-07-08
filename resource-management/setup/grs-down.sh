#!/usr/bin/env bash


set -o errexit
set -o nounset
set -o pipefail

GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

if [ -f "${GRS_ROOT}/setup/env.sh" ]; then
    source "${GRS_ROOT}/setup/env.sh"
fi

source "${GRS_ROOT}/setup/grs-util.sh"

echo "Bring down service using provider: ${CLOUD_PROVIDER}" >&2

echo "... calling verify-prereqs" >&2
verify-prereqs

echo "... calling grs-down" >&2
grs-down

echo "Done"

