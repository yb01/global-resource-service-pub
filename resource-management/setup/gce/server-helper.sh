#!/usr/bin/env bash
#
# Copyright 2022 Authors of Global Resource Service.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# create-server-instance creates the server instance. If called with
# an argument, the argument is used as the name to a reserved IP
# address for the server. (In the case of upgrade/repair, we re-use
# the same IP.)
#
# variables are set:
#   ensure-temp-dir
#   detect-project
#   get-bearer-token
function create-server-instance {
  local address=""
  [[ -n ${1:-} ]] && address="${1}"
  local internal_address=""
  [[ -n ${2:-} ]] && internal_address="${2}"

  write-server-env
  #ensure-gci-metadata-files
  create-server-instance-internal "${SERVER_NAME}" "${address}" "${internal_address}"
}

function create-server-instance-internal() {
  local gcloud="gcloud"
  local retries=5
  local sleep_sec=10

  local -r server_name="${1}"
  local -r address="${2:-}"
  local -r internal_address="${3:-}"

  local network=$(make-gcloud-network-argument \
    "${NETWORK_PROJECT}" "${REGION}" "${NETWORK}" "${SUBNETWORK:-}" \
    "${address:-}" "${ENABLE_IP_ALIASES:-}" "${IP_ALIAS_SIZE:-}")

  local metadata="server-env=${SERVICE_TEMP}/server-env.yaml"
  metadata="${metadata},user-data=${GRS_ROOT}/setup/gce/server.yaml"
  metadata="${metadata},configure-sh=${GRS_ROOT}/setup/gce/configure.sh"
  
  local disk="name=${server_name}-pd"
  disk="${disk},device-name=server-pd"
  disk="${disk},mode=rw"
  disk="${disk},boot=no"
  disk="${disk},auto-delete=no"

  for attempt in $(seq 1 ${retries}); do
    if result=$(${gcloud} compute instances create "${server_name}" \
      --project "${PROJECT}" \
      --zone "${ZONE}" \
      --machine-type "${SERVER_SIZE}" \
      --image-project="${GCE_SERVER_PROJECT}" \
      --image "${GCE_SERVER_IMAGE}" \
      --tags "${SERVER_TAG}" \
      --scopes "storage-ro,compute-rw,monitoring,logging-write" \
      --metadata-from-file "${metadata}" \
      --disk "${disk}" \
      --boot-disk-size "${SERVER_ROOT_DISK_SIZE}" \
      ${network} \
      2>&1); then
      echo "${result}" >&2

      return 0
    else
      echo "${result}" >&2
      if [[ ! "${result}" =~ "try again later" ]]; then
        echo "Failed to create server instance due to non-retryable error" >&2
        return 1
      fi
      sleep $sleep_sec
    fi
  done

  echo "Failed to create server instance despite ${retries} attempts" >&2
  return 1
}
