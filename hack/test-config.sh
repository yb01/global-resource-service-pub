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


# gcloud multiplexing for shared GCE/GKE tests.
GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

SIM_NUM=${SIM_NUM:-2}
CLIENT_NUM=${CLIENT_NUM:-2}
SERVER_NUM=${SERVER_NUM:-1}
ZONE=${GRS_GCE_ZONE:-"us-central1-a"}
REGION=${ZONE%-*}
PROJECT=${GRS_GCE_project:-"workload-controller-manager"}
INSTANCE_PREFIX=${GRS_INSTANCE_PREFIX:-grs}
SIM_INSTANCE_PREFIX=${SIM_INSTANCE_PREFIX:-"${INSTANCE_PREFIX}-sim"}
CLIENT_INSTANCE_PREFIX=${CLIENT_INSTANCE_PREFIX:-"${INSTANCE_PREFIX}-client"}
SERVER_INSTANCE_PREFIX=${SERVER_INSTANCE_PREFIX:-"${INSTANCE_PREFIX}-server"}

SOURCE_INSTANCE=${SOURCE_INSTANCE:-"sonya-grs-test-template"}
SOURCE_DISK=${SOURCE_DISK:-"sonya-grs-test-template"}
SOURCE_DISK_ZONE=${SOURCE_DISK_ZONE:-"${ZONE}"}
IMAGE_NAME=${IMAGE_NAME:-"grs-test-image"}

# simulator instance parameter
SIM_SOURCE_INSTANCE=${SIM_SOURCE_INSTANCE:-"${SOURCE_INSTANCE}"}
SIM_SOURCE_DISK=${SIM_SOURCE_DISK:-"${SOURCE_DISK}"}
SIM_SOURCE_DISK_ZONE=${SIM_SOURCE_DISK_ZONE:-"${SOURCE_DISK_ZONE}"}
SIM_IMAGE_NAME=${SIM_IMAGE_NAME:-"${IMAGE_NAME}"}
SIM_ZONE=${SIM_ZONE:-"${SIM_SOURCE_DISK_ZONE}"}
SIM_WAIT_DOWN_TIME=${SIM_WAIT_DOWN_TIME:-"5,10"}
SIM_DATA_PATTERN=${SIM_DATA_PATTERN:-"Daily"}


# client instance parameter
CLIENT_SOURCE_INSTANCE=${CLIENT_SOURCE_INSTANCE:-"${SOURCE_INSTANCE}"}
CLIENT_SOURCE_DISK=${CLIENT_SOURCE_DISK:-"${SOURCE_DISK}"}
CLIENT_SOURCE_DISK_ZONE=${CLIENT_SOURCE_DISK_ZONE:-"${SOURCE_DISK_ZONE}"}
CLIENT_IMAGE_NAME=${CLIENT_IMAGE_NAME:-"${IMAGE_NAME}"}
CLIENT_ZONE=${CLIENT_ZONE:-"${CLIENT_SOURCE_DISK_ZONE}"}

# server instance parameter
SERVER_SOURCE_INSTANCE=${SERVER_SOURCE_INSTANCE:-"sonya-grs-resourcemanagement"}
SERVER_SOURCE_DISK=${SERVER_SOURCE_DISK:-"sonya-grs-resourcemanagement"}
SERVER_SOURCE_DISK_ZONE=${SERVER_SOURCE_DISK_ZONE:-"${ZONE}"}
SERVER_IMAGE_NAME=${SERVER_IMAGE_NAME:-"grs-server-image"}
SERVER_ZONE=${SERVER_ZONE:-"${SERVER_SOURCE_DISK_ZONE}"}

#teardown parameter
SERVER_AUTO_DELETE=${SERVER_AUTO_DELETE:-true}
CLIENT_AUTO_DELETE=${CLIENT_AUTO_DELETE:-true}
SIM_AUTO_DELETE=${SIM_AUTO_DELETE:-true}
SERVERIMAGE_AUTO_DELETE=${SERVERIMAGE_AUTO_DELETE:-true}
CLIENTIMAGE_AUTO_DELETE=${CLIENTIMAGE_AUTO_DELETE:-true}
SIMIMAGE_AUTO_DELETE=${SIMIMAGE_AUTO_DELETE:-true}

#log collection parameter
DIR_ROOT=${DIR_ROOT:-"$HOME"}
SIM_LOG_DIR=${SIM_LOG_DIR:-"${DIR_ROOT}/logs"}
SERVER_LOG_DIR=${SERVER_LOG_DIR:-"${DIR_ROOT}/logs"}
CLIENT_LOG_DIR=${CLIENT_LOG_DIR:-"${DIR_ROOT}/logs"}
DES_LOG_DIR=${DES_LOG_DIR:-"${DIR_ROOT}/grs/logs/$((${SIM_NUM}*${SIM_RP_NUM}*${NODES_PER_RP}))/${SIM_DATA_PATTERN}"}
DES_LOG_INSTANCE=${DES_LOG_INSTANCE:-"sonyadev4"}
DES_LOG_INSTANCE_ZONE=${DES_LOG_INSTANCE_ZONE:-"us-central1-a"}
LOCAL_LOG_ONLY=${LOCAL_LOG_ONLY:-false}


#rune2e parameter
SIM_LOG_LEVEL=${SIM_LOG_LEVEL:-3}
SERVER_LOG_LEVEL=${SERVER_LOG_LEVEL:-3}
CLIENT_LOG_LEVEL=${CLIENT_LOG_LEVEL:-3}
SIM_CODE_ROOT=${SIM_CODE_ROOT:-"/home/sonyali/go/src/global-resource-service"}
SERVER_CODE_ROOT=${SERVER_CODE_ROOT:-"/home/sonyali/go/src/global-resource-service"}
CLIENT_CODE_ROOT=${CLIENT_CODE_ROOT:-"/home/sonyali/go/src/global-resource-service"}
SERVICE_EXTRA_ARGS=${SERVICE_EXTRA_ARGS:-}
SIM_EXTRA_ARGS=${SIM_EXTRA_ARGS:-}
CLIENT_EXTRA_ARGS=${CLIENT_EXTRA_ARGS:-}
SIM_PORT=${SIM_PORT:-"9119"}
SCHEDULER_START_DELAY=${SCHEDULER_START_DELAY:-2}   ##using between each scheduler, when starting multi scheduler on one client machine 

SIM_REGIONS=${SIM_REGIONS:-"Beijing,Shanghai"}
SIM_RP_NUM=${SIM_RP_NUM:-"10"}
NODES_PER_RP=${NODES_PER_RP:-"20000"}

SCHEDULER_REQUEST_MACHINE=${SCHEDULER_REQUEST_MACHINE:-"25000"}
SCHEDULER_REQUEST_LIMIT=${SCHEDULER_REQUEST_LIMIT:-"26000"}
SCHEDULER_NUM=${SCHEDULER_NUM:-"16"}

####if true, all service will start automaticly including resource management service, simulator, scheduler
AUTORUN_E2E=${AUTORUN_E2E:-true}
