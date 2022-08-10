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


# A library of helper functions and constant for the local config.

# Use the config file specified in $SERVICE_CONFIG_FILE, or default to
# config-default.sh.

GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..


source "${GRS_ROOT}/setup/gce/${GRS_CONFIG_FILE-"config-default.sh"}"

source "${GRS_ROOT}/setup/gce/server-helper.sh"

# These prefixes must not be prefixes of each other, so that they can be used to
# detect mutually exclusive sets of nodes.
SIMULATOR_INSTANCE_PREFIX=${NODE_INSTANCE_PREFIX:-"${INSTANCE_PREFIX}-sim"}
PROMPT_FOR_UPDATE=${PROMPT_FOR_UPDATE:-"n"}

function join_csv() {
  local IFS=','; echo "$*";
}

# This function returns the first string before the comma
function split_csv() {
  echo "$*" | cut -d',' -f1
}

# Verify prereqs
function verify-prereqs() {
  local cmd

  # we use gcloud to create the server, gsutil to stage binaries and data
  for cmd in gcloud gsutil; do
    if ! which "${cmd}" >/dev/null; then
      local resp="n"
      if [[ "${PROMPT_FOR_UPDATE}" == "y" ]]; then
        echo "Can't find ${cmd} in PATH.  Do you wish to install the Google Cloud SDK? [Y/n]"
        read resp
      fi
      if [[ "${resp}" != "n" && "${resp}" != "N" ]]; then
        curl https://sdk.cloud.google.com | bash
      fi
      if ! which "${cmd}" >/dev/null; then
        echo "Can't find ${cmd} in PATH, please fix and retry. The Google Cloud " >&2
        echo "SDK can be downloaded from https://cloud.google.com/sdk/." >&2
        exit 1
      fi
    fi
  done
  update-or-verify-gcloud
}

# Update or verify required gcloud components are installed
# at minimum required version.
# Assumed vars
#   PROMPT_FOR_UPDATE
function update-or-verify-gcloud() {
  local sudo_prefix=""
  if [ ! -w $(dirname `which gcloud`) ]; then
    sudo_prefix="sudo"
  fi
  # update and install components as needed
  if [[ "${PROMPT_FOR_UPDATE}" == "y" ]]; then
    ${sudo_prefix} gcloud ${gcloud_prompt:-} components install alpha
    ${sudo_prefix} gcloud ${gcloud_prompt:-} components install beta
    ${sudo_prefix} gcloud ${gcloud_prompt:-} components update
  else
    local version=$(gcloud version --format=json)
    python -c'
import json,sys
from distutils import version

minVersion = version.LooseVersion("1.3.0")
required = [ "alpha", "beta", "core" ]
data = json.loads(sys.argv[1])
rel = data.get("Google Cloud SDK")
if "CL @" in rel:
  print("Using dev version of gcloud: %s" %rel)
  exit(0)
if rel != "HEAD" and version.LooseVersion(rel) < minVersion:
  print("gcloud version out of date ( < %s )" % minVersion)
  exit(1)
missing = []
for c in required:
  if not data.get(c):
    missing += [c]
if missing:
  for c in missing:
    print ("missing required gcloud component \"{0}\"".format(c))
    print ("Try running `gcloud components install {0}`".format(c))
  exit(1)
    ' """${version}"""
  fi
}

# Use the gcloud defaults to find the project.  If it is already set in the
# environment then go with that.
#
# Vars set:
#   PROJECT
#   NETWORK_PROJECT
#   PROJECT_REPORTED
function detect-project() {
  if [[ -z "${PROJECT-}" ]]; then
    PROJECT=$(gcloud config list project --format 'value(core.project)')
  fi

  NETWORK_PROJECT=${NETWORK_PROJECT:-${PROJECT}}

  if [[ -z "${PROJECT-}" ]]; then
    echo "Could not detect Google Cloud Platform project.  Set the default project using " >&2
    echo "'gcloud config set project <PROJECT>'" >&2
    exit 1
  fi
  if [[ -z "${PROJECT_REPORTED-}" ]]; then
    echo "Project: ${PROJECT}" >&2
    echo "Network Project: ${NETWORK_PROJECT}" >&2
    echo "Zone: ${ZONE}" >&2
    PROJECT_REPORTED=true
  fi
}

# Example:  trap_add 'echo "in trap DEBUG"' DEBUG
# See: http://stackoverflow.com/questions/3338030/multiple-bash-traps-for-the-same-signal
function trap_add() {
  local trap_add_cmd
  trap_add_cmd=$1
  shift

  for trap_add_name in "$@"; do
    local existing_cmd
    local new_cmd

    # Grab the currently defined trap commands for this trap
    existing_cmd=$(trap -p "${trap_add_name}" |  awk -F"'" '{print $2}')

    if [[ -z "${existing_cmd}" ]]; then
      new_cmd="${trap_add_cmd}"
    else
      new_cmd="${trap_add_cmd};${existing_cmd}"
    fi

    # Assign the test. Disable the shellcheck warning telling that trap
    # commands should be single quoted to avoid evaluating them at this
    # point instead evaluating them at run time. The logic of adding new
    # commands to a single trap requires them to be evaluated right away.
    # shellcheck disable=SC2064
    trap "${new_cmd}" "${trap_add_name}"
  done
}

# Opposite of ensure-temp-dir()
cleanup-temp-dir() {
  rm -rf "${SERVICE_TEMP}"
}

# Create a temp dir that'll be deleted at the end of this bash session.
#
# Vars set:
#   SERVICE_TEMP
function ensure-temp-dir() {
  if [[ -z ${SERVICE_TEMP-} ]]; then
    SERVICE_TEMP=$(mktemp -d 2>/dev/null || mktemp -d -t grs.XXXXXX)
    trap_add cleanup-temp-dir EXIT
  fi
}

# Detect region simulators created in the instance group.
#
# Assumed vars:
#   SIM_INSTANCE_PREFIX

# Vars set:
#   SIM_NAMES
#   INSTANCE_GROUPS
function detect-sim-names() {
  detect-project
  INSTANCE_GROUPS=()
  INSTANCE_GROUPS+=($(gcloud compute instance-groups managed list \
    --project "${PROJECT}" \
    --filter "name ~ '${SIM_INSTANCE_PREFIX}-.+' AND zone:(${ZONE})" \
    --format='value(name)' || true))
  SIM_NAMES=()
  if [[ -n "${INSTANCE_GROUPS[@]:-}" ]]; then
    for group in "${INSTANCE_GROUPS[@]}"; do
      SIM_NAMES+=($(gcloud compute instance-groups managed list-instances \
        "${group}" --zone "${ZONE}" --project "${PROJECT}" \
        --format='value(instance)'))
    done
  fi

  echo "INSTANCE_GROUPS=${INSTANCE_GROUPS[*]:-}" >&2
  echo "SIM_NAMES=${SIM_NAMES[*]:-}" >&2
}

function check-network-mode() {
  local mode="$(gcloud compute networks list --filter="name=('${NETWORK}')" --project ${NETWORK_PROJECT} --format='value(x_gcloud_subnet_mode)' || true)"
  # The deprecated field uses lower case. Convert to upper case for consistency.
  echo "$(echo $mode | tr [a-z] [A-Z])"
}

function create-network() {
  if ! gcloud compute networks --project "${NETWORK_PROJECT}" describe "${NETWORK}" &>/dev/null; then
    # The network needs to be created synchronously or we have a race. The
    # firewalls can be added concurrent with instance creation.
    local network_mode="auto"
    if [[ "${CREATE_CUSTOM_NETWORK:-}" == "true" ]]; then
      network_mode="custom"
    fi
    echo "Creating new ${network_mode} network: ${NETWORK}"
    gcloud compute networks create --project "${NETWORK_PROJECT}" "${NETWORK}" --subnet-mode="${network_mode}"
  else
    PREEXISTING_NETWORK=true
    PREEXISTING_NETWORK_MODE="$(check-network-mode)"
    echo "Found existing network ${NETWORK} in ${PREEXISTING_NETWORK_MODE} mode."
  fi
}

function create-subnetworks() {
  case ${ENABLE_IP_ALIASES} in
    true) echo "IP aliases are enabled. Creating subnetworks.";;
    false)
      echo "IP aliases are disabled."
      if [[ "${ENABLE_BIG_CLUSTER_SUBNETS}" = "true" ]]; then
        if [[  "${PREEXISTING_NETWORK}" != "true" ]]; then
          expand-default-subnetwork
        else
          echo "${color_yellow}Using pre-existing network ${NETWORK}, subnets won't be expanded to /19!${color_norm}"
        fi
      elif [[ "${CREATE_CUSTOM_NETWORK:-}" == "true" && "${PREEXISTING_NETWORK}" != "true" ]]; then
          gcloud compute networks subnets create "${SUBNETWORK}" --project "${NETWORK_PROJECT}" --region "${REGION}" --network "${NETWORK}" --range "${NODE_IP_RANGE}"
      fi
      return;;
    *) echo "${color_red}Invalid argument to ENABLE_IP_ALIASES${color_norm}"
       exit 1;;
  esac

  # Look for the alias subnet, it must exist and have a secondary
  # range configured.
  local subnet=$(gcloud compute networks subnets describe \
    --project "${NETWORK_PROJECT}" \
    --region ${REGION} \
    ${IP_ALIAS_SUBNETWORK} 2>/dev/null)
  if [[ -z ${subnet} ]]; then
    echo "Creating subnet ${NETWORK}:${IP_ALIAS_SUBNETWORK}"
    gcloud compute networks subnets create \
      ${IP_ALIAS_SUBNETWORK} \
      --description "Automatically generated subnet for ${INSTANCE_PREFIX} cluster. This will be removed on cluster teardown." \
      --project "${NETWORK_PROJECT}" \
      --network ${NETWORK} \
      --region ${REGION} \
      --range ${NODE_IP_RANGE} \
      --secondary-range "pods-default=${CLUSTER_IP_RANGE}" \
      --secondary-range "services-default=${SERVICE_CLUSTER_IP_RANGE}"
    echo "Created subnetwork ${IP_ALIAS_SUBNETWORK}"
  else
    if ! echo ${subnet} | grep --quiet secondaryIpRanges; then
      echo "${color_red}Subnet ${IP_ALIAS_SUBNETWORK} does not have a secondary range${color_norm}"
      exit 1
    fi
  fi
}

# Robustly try to create a static ip.
# $1: The name of the ip to create
# $2: The name of the region to create the ip in.
function create-static-ip() {
  detect-project
  local attempt=0
  local REGION="$2"
  while true; do
    if gcloud compute addresses create "$1" \
      --project "${PROJECT}" \
      --region "${REGION}" -q > /dev/null; then
      # successful operation - wait until it's visible
      start="$(date +%s)"
      while true; do
        now="$(date +%s)"
        # Timeout set to 15 minutes
        if [[ $((now - start)) -gt 900 ]]; then
          echo "Timeout while waiting for server IP visibility"
          exit 2
        fi
        if gcloud compute addresses describe "$1" --project "${PROJECT}" --region "${REGION}" >/dev/null 2>&1; then
          break
        fi
        echo "server IP not visible yet. Waiting..."
        sleep 5
      done
      break
    fi

    if gcloud compute addresses describe "$1" \
      --project "${PROJECT}" \
      --region "${REGION}" >/dev/null 2>&1; then
      # it exists - postcondition satisfied
      break
    fi

    if (( attempt > 4 )); then
      echo -e "${color_red}Failed to create static ip $1 ${color_norm}" >&2
      exit 2
    fi
    attempt=$(($attempt+1))
    echo -e "${color_yellow}Attempt $attempt failed to create static ip $1. Retrying.${color_norm}" >&2
    sleep $(($attempt * 5))
  done
}

# Instantiate resource management service
#

function grs-up() {
  ensure-temp-dir
  detect-project
  create-network
  create-resourcemanagement-server
  create-region-simulator
}

# tear done resource management service

function grs-down() {
  detect-project

  echo "Bringing down resource management service"
  set +e  # Do not stop on error
  
  # Get the name of the managed instance group template and delete
  local templates=$(get-template "${PROJECT}")

  local all_instance_groups=(${INSTANCE_GROUPS[@]:-})
  for group in ${all_instance_groups[@]:-}; do
    {
      if gcloud compute instance-groups managed describe "${group}" --project "${PROJECT}" --zone "${ZONE}" &>/dev/null; then
        gcloud compute instance-groups managed delete \
          --project "${PROJECT}" \
          --quiet \
          --zone "${ZONE}" \
          "${group}"
      fi
    } &
  done

  wait-for-jobs || {
      echo -e "Failed to delete instance template(s)." >&2
    }


  # Check if this are any remaining server replicas.
  local REMAINING_SERVER_COUNT=0
#: <<'EOF'
  REMAINING_SERVER_COUNT=$(gcloud compute instances list \
    --project "${PROJECT}" \
    --filter="name ~ '$(get-replica-name-regexp)'" \
    --format "value(zone)" | wc -l)

  if [[ "${REMAINING_SERVER_COUNT}" -ge 1 ]]; then
    local instance_names=$(get-all-replica-names)

    for instance_name in ${instance_names[@]:-}; do
    {
      if gcloud compute instances describe "${instance_name}" --zone "${ZONE}" --project "${PROJECT}" &>/dev/null; then
        gcloud compute instances delete \
          --project "${PROJECT}" \
          --zone "${ZONE}" \
          --quiet \
          "${instance_name}"
      fi
    }
    done
    
    wait-for-jobs || {
      echo -e "Failed to delete server(s)." >&2
    }
  fi
#EOF

  REMAINING_SERVER_COUNT=$(gcloud compute instances list \
    --project "${PROJECT}" \
    --filter="name ~ '$(get-replica-name-regexp)'" \
    --format "value(zone)" | wc -l)

  if [[ "${REMAINING_SERVER_COUNT}" -eq 0 ]]; then    
    # Delete the server's reserved IP
    if gcloud compute addresses describe "${SERVER_NAME}-ip" --region "${REGION}" --project "${PROJECT}" &>/dev/null; then
      echo "Deleting the server's reserved IP"
      gcloud compute addresses delete \
        --project "${PROJECT}" \
        --region "${REGION}" \
        --quiet \
        "${SERVER_NAME}-ip"
    fi

    # Delete the server's pd
    if gcloud compute disks describe "${SERVER_NAME}-pd" --zone "${ZONE}" --project "${PROJECT}" &>/dev/null; then
      echo "Deleting the server's pd"
      gcloud compute disks delete \
        --project "${PROJECT}" \
        --zone "${ZONE}" \
        --quiet \
        "${SERVER_NAME}-pd"
    fi
  fi

  set -e
}

function get-replica-name-regexp() {
  echo "^${SERVER_NAME}(-...)?"
}

function get-all-replica-names() {
  echo $(gcloud compute instances list \
    --project "${PROJECT}" \
    --filter="name ~ '$(get-replica-name-regexp)'" \
    --format "value(name)" | tr "\n" "," | sed 's/,$//')
}

# Gets the instance templates in use by the service. It echos the template names
# so that the function output can be used.

function get-template() {
  local linux_filter="${SIM_INSTANCE_PREFIX}-(extra-)?template(-)?"
  
  gcloud compute instance-templates list \
    --filter="name ~ '${linux_filter}'" \
    --project="${1}" --format='value(name)'
}

function create-resourcemanagement-server() {
  echo "Starting rersource management server"
  
  # We have to make sure the disk is created before creating the server VM, so
  # run this in the foreground.
  gcloud compute disks create "${SERVER_NAME}-pd" \
    --project "${PROJECT}" \
    --zone "${ZONE}" \
    --type "${SERVER_DISK_TYPE}" \
    --size "${SERVER_DISK_SIZE}"

  # Reserve the server's IP so that it can later be transferred to another VM
  create-static-ip "${SERVER_NAME}-ip" "${REGION}"
  SERVER_RESERVED_IP=$(gcloud compute addresses describe "${SERVER_NAME}-ip" \
    --project "${PROJECT}" --region "${REGION}" -q --format='value(address)')

  create-server-instance "${SERVER_RESERVED_IP}" 


}

function create-region-simulator() {
  echo "Starting region simulatotrs"
    #create-nodes-template
    #create-linux-nodes
}

# Quote something appropriate for a yaml string.
#
# TODO(zmerlynn): Note that this function doesn't so much "quote" as
# "strip out quotes", and we really should be using a YAML library for
# this, but PyYAML isn't shipped by default, and *rant rant rant ... SIGH*
function yaml-quote {
  echo "'$(echo "${@:-}" | sed -e "s/'/''/g")'"
}

function write-server-env {
  build-server-env "server" "${SERVICE_TEMP}/server-env.yaml"
}

function write-sim-env {
  build-server-env "sim" "${SERVICE_TEMP}/server-env.yaml"
}

function build-server-env {
  local server="$1"
  local file="$2"

  rm -f ${file}
  cat >$file <<EOF
INSTANCE_PREFIX:  $(yaml-quote ${INSTANCE_PREFIX:-})
GOLANG_VERSION:   $(yaml-quote ${GOLANG_VERSION:-})
REDIS_VERSION:    $(yaml-quote ${REDIS_VERSION:-})
EOF

  if [[ "${server}" == "server" ]]; then
    cat >>$file <<EOF
SERVER_NAME: $(yaml-quote ${SERVER_NAME})
EOF
  fi

  if [[ "${server}" == "sim" ]]; then
    cat >>$file <<EOF
SIM_INSTANCE_PREFIX:  $(yaml-quote ${SIM_INSTANCE_PREFIX})
EOF
  fi
}


# Format the string argument for gcloud network.
function make-gcloud-network-argument() {
  local network_project="$1"
  local region="$2"
  local network="$3"
  local subnet="$4"
  local address="$5"          # optional
  local enable_ip_alias="$6"  # optional
  local alias_size="$7"       # optional

  local networkURL="projects/${network_project}/global/networks/${network}"
  local subnetURL="projects/${network_project}/regions/${region}/subnetworks/${subnet:-}"

  local ret=""

  if [[ "${enable_ip_alias}" == 'true' ]]; then
    ret="--network-interface"
    ret="${ret} network=${networkURL}"
    if [[ "${address:-}" == "no-address" ]]; then
      ret="${ret},no-address"
    else
      ret="${ret},address=${address:-}"
    fi
    ret="${ret},subnet=${subnetURL}"
    ret="${ret},aliases=pods-default:${alias_size}"
    ret="${ret} --no-can-ip-forward"
  else
    if [[ -n ${subnet:-} ]]; then
      ret="${ret} --subnet ${subnetURL}"
    else
      ret="${ret} --network ${networkURL}"
    fi

    ret="${ret} --can-ip-forward"
    if [[ -n ${address:-} ]] && [[ "$address" != "no-address" ]]; then
      ret="${ret} --address ${address}"
    fi
  fi

  echo "${ret}"
}

# Wait for background jobs to finish. Return with
# an error status if any of the jobs failed.
function wait-for-jobs() {
  local fail=0
  local job
  for job in $(jobs -p); do
    wait "${job}" || fail=$((fail + 1))
  done
  return ${fail}
}