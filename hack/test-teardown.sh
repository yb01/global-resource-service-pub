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


### Only support gcloud 
### Please ensure gcloud is installed before run this script
GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${GRS_ROOT}/hack/test-config.sh"

echo "Tear down server, region simulator and clients... "

function delete-image {
        local image_name="$1"
        if gcloud compute images describe "${image_name}" --project "${PROJECT}" &>/dev/null; then
                echo "Deleting existing template image: ${image_name}"
                gcloud compute images delete \
                        "${image_name}" \
                        --project "${PROJECT}" \
                        --quiet 
        fi
}

function delete-instance-template {
        local template_name="$1"
        if gcloud compute instance-templates describe "${template_name}" --project "${PROJECT}" &>/dev/null; then
                echo "Try to Delete existing instance template: ${template_name}, if failed, Please manully delete instance group first!"
                gcloud compute instance-templates delete \
                        "${template_name}" \
                        --project "${PROJECT}" \
                        --quiet 
        fi
}

function delete-instance-groups {
        local group_name="$1"
        local zone="$2"
        if gcloud compute instance-groups managed describe "${group_name}" --project "${PROJECT}" --zone "${zone}" &>/dev/null; then
                 echo "Try to Delete existing instance groups: ${group_name}"
                 gcloud compute instance-groups managed delete \
                         "${group_name}" \
                         --project "${PROJECT}" \
                         --zone "${zone}" \
                         --quiet 
        fi
}

function delete-vm-instance {
        local vm_name="$1"
        local instance_zone="$2"
        if gcloud compute instances describe "${vm_name}" --project "${PROJECT}" --zone "${instance_zone}" &>/dev/null; then
                echo "Try to delete existing instance : ${vm_name}"
                gcloud compute instances delete \
                        "${vm_name}" \
                        --project "${PROJECT}" \
                        --zone "${instance_zone}"  \
                        --quiet 
        fi
}    

function delete-machine-image {
        local image_name="$1"
        if gcloud compute machine-images describe "${image_name}" --project "${PROJECT}" &>/dev/null; then
                echo "Deleting existing machine image: ${image_name}"
                gcloud compute machine-images delete "${image_name}" \
                        --project "${PROJECT}" \
                        --quiet
        fi
}

###############
#   main function
###############

if [[ "${SIM_AUTO_DELETE}" == "true" && ${SIM_NUM} -gt 0 ]]; then
        echo "Deleting simulator resources"
        IFS=','; INSTANCE_SIM_ZONE=($SIM_ZONE); unset IFS;
        index=0
        for zone in "${INSTANCE_SIM_ZONE[@]}"; do
                delete-vm-instance "${SIM_INSTANCE_PREFIX}-${zone}-${index}" "${zone}" &
                index=$(($index + 1)) 
        done
        if [ "${SIMIMAGE_AUTO_DELETE}" == "true" ]; then
                #waiting 60 seconds to get all instances deleted before delete images
                sleep 60
                delete-machine-image  "${SIM_IMAGE_NAME}"
        fi
fi

if [[ "${CLIENT_AUTO_DELETE}" == "true"  && ${CLIENT_NUM} -gt 0 ]]; then
        echo "Deleting client scheduler resources"
        IFS=','; INSTANCE_CLIENT_ZONE=($CLIENT_ZONE); unset IFS;
        index=0
        for zone in "${INSTANCE_CLIENT_ZONE[@]}"; do
                delete-vm-instance "${CLIENT_INSTANCE_PREFIX}-${zone}-${index}" "${zone}" &
                index=$(($index + 1)) 
        done
        if [ "${CLIENTIMAGE_AUTO_DELETE}" == "true" ]; then
                #waiting 60 seconds to get all instances deleted before delete images
                sleep 60
                delete-machine-image  "${CLIENT_IMAGE_NAME}"
        fi
fi

if [[ "${SERVER_AUTO_DELETE}" == "true"  && ${SERVER_NUM} -gt 0 ]]; then
        echo "Deleting server resources"
        IFS=','; INSTANCE_SERVER_ZONE=($SERVER_ZONE); unset IFS;
        index=0
        for zone in "${INSTANCE_SERVER_ZONE[@]}"; do
                delete-vm-instance "${SERVER_INSTANCE_PREFIX}-${zone}-${index}" "${zone}" &
                index=$(($index + 1)) 
        done
        if [ "${SERVERIMAGE_AUTO_DELETE}" == "true" ]; then
                #waiting 60 seconds to get all instances deleted before delete images
                sleep 60
                delete-machine-image  "${SERVER_IMAGE_NAME}"
        fi
fi

echo "Done. All resources deleted successfully"

