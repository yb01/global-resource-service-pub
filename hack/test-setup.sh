#!/usr/bin/env bash

### Only support gcloud 
### Please ensure gcloud is installed before run this script

GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")

source "${GRS_ROOT}/test-config.sh"

function create-image {
        local image_name="$1"
        local source_disk="$2"
        local source_disk_zone="$3"
        echo "Check and create images  with image_name: ${image_name}, source_disk: ${source_disk}."
        if gcloud compute images describe "${image_name}" --project "${PROJECT}" &>/dev/null; then
                echo "Image name: ${image_name} already exist, using existing one."
        else
                gcloud compute images create \
                        "${image_name}" \
                        --project "${PROJECT}" \
                        --source-disk "${source_disk}" \
                        --source-disk-zone "${source_disk_zone}" \
                        --quiet
        fi
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
        echo "Check and create instance groups"
        if gcloud compute instance-groups managed describe "${group_name}" --project "${PROJECT}" &>/dev/null; then
                echo "Try to Delete existing instance group: ${group_name}"
                gcloud compute instance-groups managed delete \
                        "${group_name}" \
                        --project "${PROJECT}" \
                        --quiet 
        fi
        gcloud compute instance-groups managed create \
                "${group_name}" \
                --project "${PROJECT}" \
                --zone "${ZONE}" \
                --template "${template_name}" \
                --size "${instance_num}" \
                --quiet
}

function attach-ipaddress {
        local group_name="$1"
        echo "Attaching reserved ip address to instances"
        instance_names=()
        instance_names=($(gcloud compute instance-groups managed list-instances \
        "${group_name}" --zone "${ZONE}" --project "${PROJECT}" \
        --format='value(instance)'))

        ip_index=0
        for name in "${instance_names[@]}"; do
                echo $name
                gcloud compute instances delete-access-config \
                $name \
                --zone ${ZONE} \
                --project "${PROJECT}" \
                --access-config-name="External NAT"

                gcloud compute instances add-access-config \
                $name \
                --zone ${ZONE} \
                --project "${PROJECT}" \
                --access-config-name "External NAT" \
                --address "${IP_ADDRESS[$ip_index]}"
                echo "$ip_index,${IP_ADDRESS[$ip_index]}"
                ip_index=$(($ip_index + 1))
        done
}


###############
#   main function
###############

#verify reserved ip address number is greater then simulator vm number
IFS=','; IP_ADDRESS=($SIM_RESERVED_IP); unset IFS;

if [ ${#IP_ADDRESS[@]} -lt ${SIM_NUM} ]; then
        echo "Reserved IP address number less then simulator vm number, Please double check."
        exit 1
fi

if [ ${SIM_NUM} -gt 0 ]; then
        echo "starting region simulator... "
        create-image "${SIM_IMAGE_NAME}" "${SIM_SOURCE_DISK}" "${SIM_SOURCE_DISK_ZONE}"
        create-template "${SIM_INSTANCE_PREFIX}-template" "${SIM_SOURCE_INSTANCE}" "${SIM_SOURCE_DISK}" "${SIM_IMAGE_NAME}" "${SIM_SOURCE_DISK_ZONE}"
        create-instance-group "${SIM_INSTANCE_PREFIX}-instance-group" "${SIM_INSTANCE_PREFIX}-template" "${SIM_NUM}"
        echo "attach reserved IP to region simulator"
        attach-ipaddress "${SIM_INSTANCE_PREFIX}-instance-group"
fi

if [ ${CLIENT_NUM} -gt 0 ]; then
        echo "starting client scheduler... "
        create-image "${CLIENT_IMAGE_NAME}" "${CLIENT_SOURCE_DISK}" "${CLIENT_SOURCE_DISK_ZONE}"
        create-template "${CLIENT_INSTANCE_PREFIX}-template" "${CLIENT_SOURCE_INSTANCE}" "${CLIENT_SOURCE_DISK}" "${CLIENT_IMAGE_NAME}" "${CLIENT_SOURCE_DISK_ZONE}"
        create-instance-group "${CLIENT_INSTANCE_PREFIX}-instance-group" "${CLIENT_INSTANCE_PREFIX}-template" "${CLIENT_NUM}"
fi

if [ ${SERVER_NUM} -gt 0 ]; then
        echo "starting resource management service... "
        create-image "${SERVER_IMAGE_NAME}" "${SERVER_SOURCE_DISK}" "${SERVER_SOURCE_DISK_ZONE}"
        create-template "${SERVER_INSTANCE_PREFIX}-template" "${SERVER_SOURCE_INSTANCE}" "${SERVER_SOURCE_DISK}" "${SERVER_IMAGE_NAME}" "${SERVER_SOURCE_DISK_ZONE}"
        create-instance-group "${SERVER_INSTANCE_PREFIX}-instance-group" "${SERVER_INSTANCE_PREFIX}-template" "${SERVER_NUM}"
fi

