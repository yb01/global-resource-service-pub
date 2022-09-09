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

### Instructions
### Please double check if these env are correct before run this script
### Suggested SERVER_ZONE: us-central1-a
### Suggested SIM_ZONE, CLIENT_ZONE: us-central1-a,us-west1-c,us-west2-a,us-west4-a,us-west3-c,us-east1-b,us-east4-b,asia-east1-b,asia-south2-a,australia-southeast1-b,northamerica-northeast1-a...
###############
### export SIM_NUM=5 CLIENT_NUM=6 SERVER_NUM=1 
### export SERVER_ZONE=us-central1-a   
### export SIM_ZONE=us-central1-a,us-west1-c,us-west2-a,us-west4-a,us-west3-c
### export CLIENT_ZONE=us-west3-b,us-east4-b,asia-south2-a,australia-southeast1-b,northamerica-northeast1-a,us-west2-c
################

GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${GRS_ROOT}/hack/test-config.sh"
source "${GRS_ROOT}/hack/lib/util.sh"

function create-image {
        local image_name="$1"
        local source_disk="$2"
        local source_disk_zone="$3"
        echo "Check and create template images with image_name: ${image_name}, source_disk: ${source_disk}."
        if gcloud compute images describe "${image_name}" --project "${PROJECT}" &>/dev/null; then
                echo "Image name: ${image_name} already exist, using existing one."
        else
                gcloud compute images create \
                        "${image_name}" \
                        --project "${PROJECT}" \
                        --source-disk "${source_disk}" \
                        --source-disk-zone "${source_disk_zone}" \
                        --storage-location "us" \
                        --force \
                        --quiet
        fi
}

function create-machine-image {
        local image_name="$1"
        local source_instance="$2"
        local source_instance_zone="$3"
        echo "Check and create machine images  with image_name: ${image_name}, source_instance: ${source_instance}."
        if gcloud compute machine-images describe "${image_name}" --project "${PROJECT}" &>/dev/null; then
                echo "Image name: ${image_name} already exist, using existing one."
        else
                gcloud compute machine-images create "${image_name}" \
                        --source-instance "${source_instance}" \
                        --source-instance-zone "${source_instance_zone}" \
                        --project "${PROJECT}" \
                        --quiet
        fi
}

function create-vm-instance {
        local vm_name="$1"
        local instance_zone="$2"
        local machine_image="$3"
        if gcloud compute instances describe "${vm_name}" --project "${PROJECT}" --zone "${instance_zone}" &>/dev/null; then
                echo "Try to Delete existing instance : ${vm_name}"
                gcloud compute instance delete \
                        "${vm_name}" \
                        --project "${PROJECT}" \
                        --zone "${instance_zone}"  \
                        --quiet 
        fi

        gcloud beta compute instances create "${vm_name}" \
                --zone "${instance_zone}" \
                --project "${PROJECT}" \
                --source-machine-image="${machine_image}" \
                --quiet
}

function create-template {
        local template_name="$1"
        local source_instance="$2"
        local source_disk="$3"
        local image_name="$4"
        local source_instance_zone="$5"
        echo "Check and create instance templates"
        if gcloud compute instance-templates describe "${template_name}" --project "${PROJECT}" &>/dev/null; then
                gcloud compute instance-templates delete \
                "${template_name}" \
                --project "${PROJECT}" \
                --quiet 
        fi
        gcloud compute instance-templates create \
                "${template_name}" \
                --project "${PROJECT}" \
                --source-instance "${source_instance}" \
                --source-instance-zone "${source_instance_zone}" \
                --configure-disk=device-name="${source_disk}",instantiate-from=custom-image,custom-image="projects/${PROJECT}/global/images/${image_name}" \
                --quiet
        
}

function create-instance-group {
        local group_name="$1"
        local template_name="$2"
        local instance_num="$3"
        local instance_zone="$4"
        echo "Check and create instance group"
        if gcloud compute instance-groups managed describe "${group_name}" --project "${PROJECT}" &>/dev/null; then
                echo "Try to Delete existing instance group: ${group_name}"
                gcloud compute instance-groups managed delete \
                        "${group_name}" \
                        --project "${PROJECT}" \
                        --quiet 
        fi
        gcloud beta compute instance-groups managed create \
                "${group_name}" \
                --project "${PROJECT}" \
                --zone "${instance_zone}" \
                --template "${template_name}" \
                --size "${instance_num}" \
                --stateful-internal-ip enabled,auto-delete=on-permanent-instance-deletion \
                --stateful-external-ip enabled,auto-delete=on-permanent-instance-deletion \
                --quiet
}

function start-instance-redis {
        local name="$1"
        local zone="$2"
        cmd="sudo systemctl restart redis-server.service"
        ssh-config "${cmd}" "${name}" "${zone}"
        cmd="redis-cli flushall"
        ssh-config "${cmd}" "${name}" "${zone}"
}

function start-mig-redis {
        local group_name="$1"
        local zone="$2"
        instance_names=()
        instance_names=($(gcloud compute instance-groups managed list-instances \
        "${group_name}" --zone "${zone}" --project "${PROJECT}" \
        --format='value(instance)'))
        for name in "${instance_names[@]}"; do
                start-instance-redis "${name}" "${zone}"
        done
}

function get-instance-ip {
        local name="$1"
        local zone="$2"

        instance_ip=($(gcloud compute instances describe "${name}" \
                --zone "${zone}" \
                --format='get(networkInterfaces[0].accessConfigs[0].natIP)' \
                --quiet))
        echo "${instance_ip}"
}

function get-mig-ips {
        local group_name="$1"
        local zone="$2"

        instance_names=()
        instance_names=($(gcloud compute instance-groups managed list-instances \
        "${group_name}" --zone "${zone}" --project "${PROJECT}" \
        --format='value(instance)'))

        mig_ips=""
        for name in "${instance_names[@]}"; do
                mig_ips+="$(get-instance-ip ${name} ${zone}),"
        done
        echo "${mig_ips}"
}

function get-mig-urls {
        local group_name="$1"
        local zone="$2"

        instance_names=()
        instance_names=($(gcloud compute instance-groups managed list-instances \
        "${group_name}" --zone "${zone}" --project "${PROJECT}" \
        --format='value(instance)'))

        mig_urls=""
        for name in "${instance_names[@]}"; do
                mig_urls+="$(get-instance-ip ${name} ${zone}):${SIM_PORT},"
        done
        echo "${mig_urls}"
}

function get-mig-names {
        local group_name="$1"
        local zone="$2"

        instance_names=()
        instance_names=($(gcloud compute instance-groups managed list-instances \
        "${group_name}" --zone "${zone}" --project "${PROJECT}" \
        --format='value(instance)'))

        mig_names=""
        for name in "${instance_names[@]}"; do
                mig_names+="${name},"
        done
        echo "${mig_names}"
}

###############
#   main function
###############

IFS=','; INSTANCE_SERVER_ZONE=($SERVER_ZONE); unset IFS;

if [ ${#INSTANCE_SERVER_ZONE[@]} != ${SERVER_NUM} ]; then
        echo "Server zone must be same as server number, Please double check."
        exit 1
fi

IFS=','; INSTANCE_SIM_ZONE=($SIM_ZONE); unset IFS;
IFS=','; SIM_DOWN_TIME_LIST=($SIM_WAIT_DOWN_TIME); unset IFS;

if [ ${#INSTANCE_SIM_ZONE[@]} != ${SIM_NUM} ]; then
        echo "Simulator zone must be same as Simulator number, Please double check."
        exit 1
fi

if [[ "${SIM_DATA_PATTERN}" == "Outage" && "${#SIM_DOWN_TIME_LIST[@]}" != "${SIM_NUM}" ]]; then
        echo "The number of simulator wait time for make rp down:SIM_WAIT_DOWN_TIME must be same as Simulator number, Please double check."
        exit 1
fi

IFS=','; INSTANCE_CLIENT_ZONE=($CLIENT_ZONE); unset IFS;

if [ ${#INSTANCE_CLIENT_ZONE[@]} != ${CLIENT_NUM} ]; then
        echo "Client zone must be same as client number, Please double check."
        exit 1
fi

if [ ${SIM_NUM} -gt 0 ]; then
        echo "starting region simulator... "
        create-machine-image "${SIM_IMAGE_NAME}" "${SIM_SOURCE_INSTANCE}" "${SIM_SOURCE_DISK_ZONE}"
        index=0
        for zone in "${INSTANCE_SIM_ZONE[@]}"; do
                create-vm-instance "${SIM_INSTANCE_PREFIX}-${zone}-${index}" "${zone}" "${SIM_IMAGE_NAME}" &
                index=$(($index + 1)) 
        done

        echo "waiting 10 seconds to get all simulator resources created"
        sleep 10
fi

if [ ${CLIENT_NUM} -gt 0 ]; then
        echo "starting client scheduler... "
        create-machine-image "${CLIENT_IMAGE_NAME}" "${CLIENT_SOURCE_INSTANCE}" "${CLIENT_SOURCE_DISK_ZONE}"
        index=0
        for zone in "${INSTANCE_CLIENT_ZONE[@]}"; do
                create-vm-instance "${CLIENT_INSTANCE_PREFIX}-${zone}-${index}" "${zone}" "${CLIENT_IMAGE_NAME}" &
                index=$(($index + 1)) 
        done

        echo "waiting 10 seconds to get all client scheduler resources created"
        sleep 10
fi

if [ ${SERVER_NUM} -gt 0 ]; then
        echo "starting resource management service... "
        create-machine-image "${SERVER_IMAGE_NAME}" "${SERVER_SOURCE_INSTANCE}" "${SERVER_SOURCE_DISK_ZONE}"
        index=0
        for zone in "${INSTANCE_SERVER_ZONE[@]}"; do
                create-vm-instance "${SERVER_INSTANCE_PREFIX}-${zone}-${index}" "${zone}" "${SERVER_IMAGE_NAME}" &
                index=$(($index + 1)) 
        done
        echo "waiting 10 seconds to get all server resources created"
        sleep 10
fi

 
echo "Waiting 90 seconds to get all resource started"
sleep 90

RESOURCE_URLS=""
MASTER_IP=""
SERVICE_URL=""
SERVER_ZONE=""
if [ ${SIM_NUM} -gt 0 ]; then
        SIM_IPS=""
        RESOURCE_URLS=""
        index=0
        for zone in "${INSTANCE_SIM_ZONE[@]}"; do
                SIM_IPS+="$(get-instance-ip ${SIM_INSTANCE_PREFIX}-${zone}-${index} ${zone}),"
                RESOURCE_URLS+="$(get-instance-ip ${SIM_INSTANCE_PREFIX}-${zone}-${index} ${zone}):${SIM_PORT},"
                index=$((index + 1)) 
        done

        echo "Simulators started at ip addresss: ${SIM_IPS%,}"
fi

if [ ${CLIENT_NUM} -gt 0 ]; then
        CLIENT_IPS=""
        index=0
        for zone in "${INSTANCE_CLIENT_ZONE[@]}"; do
                CLIENT_IPS+="$(get-instance-ip ${CLIENT_INSTANCE_PREFIX}-${zone}-${index} ${zone}),"
                index=$((index + 1)) 
        done
        echo "Client schedulers started at ip addresss: ${CLIENT_IPS%,}"
fi

if [ ${SERVER_NUM} -gt 0 ]; then
        SERVER_IPS=""
        SERVER_NAMES=""
        SERVICE_ZONE="${INSTANCE_SERVER_ZONE[0]}"
        index=0
        for zone in "${INSTANCE_SERVER_ZONE[@]}"; do
                start-instance-redis "${SERVER_INSTANCE_PREFIX}-${zone}-${index}" "${zone}"
                SERVER_IPS+="$(get-instance-ip ${SERVER_INSTANCE_PREFIX}-${zone}-${index} ${zone}),"
                SERVER_NAMES+="${SERVER_INSTANCE_PREFIX}-${zone}-${index},"
                index=$((index + 1)) 
        done
        echo "Servers started at ip addresss: ${SERVER_IPS%,}"
fi

####Most cloud doesn't support binding to public IP, so using machine name to listen and bind service for now
export MASTER_IP="${SERVER_NAMES%%,*}"
export SERVICE_URL="${SERVER_IPS%%,*}:8080"
export RESOURCE_URLS="${RESOURCE_URLS%,}"
export SERVICE_ZONE                     # the zone on resosurce management serevice instance

echo "Done to create and start all required resouce"

if [ "${AUTORUN_E2E}" == "true" ]; then
        #Starting e2e testing
        echo "Starting service using args: --master_ip=${MASTER_IP} --resource_urls=${RESOURCE_URLS}"
        echo "Starting scheduler using args: --service_url=${SERVICE_URL}"
        "${GRS_ROOT}/hack/test-rune2e.sh"
else
        echo "You can start service using args: --master_ip=${MASTER_IP} --resource_urls=${RESOURCE_URLS}"
        echo "You can start scheduler using args: --service_url=${SERVICE_URL}"
fi


