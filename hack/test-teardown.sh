#!/usr/bin/env bash

### Only support gcloud 
### Please ensure gcloud is installed before run this script
GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")

source "${GRS_ROOT}/test-config.sh"

echo "${GRS_ROOT}/test-config.sh"

echo "Tear down region simulator and clients... "

function delete-image {
        local image_name="$1"
        if gcloud compute images describe "${image_name}" --project "${PROJECT}" &>/dev/null; then
                echo "Deleting existing instance image: ${image_name}"
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
        if gcloud compute instance-groups managed describe "${group_name}" --project "${PROJECT}" --zone "${ZONE}" &>/dev/null; then
                echo "Try to Delete existing instance groups: ${group_name}"
                gcloud compute instance-groups managed delete \
                        "${group_name}" \
                        --project "${PROJECT}" \
                        --zone "${ZONE}" \
                        --quiet 
        fi
}



###############
#   main function
###############

if [ "${SIM_AUTO_DELETE}" == "true" ]; then
        echo "Deleting simulator resources"
        delete-instance-groups "${SIM_INSTANCE_PREFIX}-instance-group"
        delete-instance-template "${SIM_INSTANCE_PREFIX}-template"
        if [ "${SIMIMAGE_AUTO_DELETE}" == "true" ]; then
                delete-image "${SIM_IMAGE_NAME}"
        fi
fi

if [ "${CLIENT_AUTO_DELETE}" == "true" ]; then
        echo "Deleting client scheduler resources"
        delete-instance-groups "${CLIENT_INSTANCE_PREFIX}-instance-group"
        delete-instance-template "${CLIENT_INSTANCE_PREFIX}-template"
        if [ "${CLIENTIMAGE_AUTO_DELETE}" == "true" ]; then
                delete-image "${CLIENT_IMAGE_NAME}"
        fi
fi

if [ "${SERVER_AUTO_DELETE}" == "true" ]; then
        echo "Deleting server resources"
        delete-instance-groups "${SERVER_INSTANCE_PREFIX}-instance-group"
        delete-instance-template "${SERVER_INSTANCE_PREFIX}-template"
        if [ "${SERVERIMAGE_AUTO_DELETE}" == "true" ]; then
                delete-image "${SERVER_IMAGE_NAME}"
        fi
fi

echo "Done. All resources deleted successfully"

