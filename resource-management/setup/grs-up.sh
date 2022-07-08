#!/usr/bin/env bash


set -o errexit
set -o nounset
set -o pipefail

GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

if [ -f "${GRS_ROOT}/setup/env.sh" ]; then
    source "${GRS_ROOT}/setup/env.sh"
fi

source "${GRS_ROOT}/setup/grs-util.sh"

if [ -z "${ZONE-}" ]; then
  echo "... Starting cluster using provider: ${CLOUD_PROVIDER}" >&2
else
  echo "... Starting cluster in ${ZONE} using provider ${CLOUD_PROVIDER}" >&2
fi

echo "... calling verify-prereqs" >&2
verify-prereqs

echo "... calling grs-up" >&2
grs-up

echo -e "Done, resource management service is running!\n" >&2

echo

exit 0
