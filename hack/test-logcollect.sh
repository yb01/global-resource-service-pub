#!/usr/bin/env bash

### Only support gcloud 
### Please ensure gcloud is installed before run this script
GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${GRS_ROOT}/hack/test-config.sh"
source "${GRS_ROOT}/hack/lib/util.sh"

echo "Collect test logs from server,region simulator and clients... "

function collect-log-instance {
    local source_name="$1"
    local source_zone="$2"
    local source_dir="$3"
    local destination_dir="$4"
    
    if [[ ! -e "${destination_dir}" ]]; then
        mkdir -p $destination_dir
    elif [[ ! -d $destination_dir ]]; then
        echo "$destination_dir already exists but is not a directory" 1>&2
        exit
    fi
    
    gcloud compute scp --zone "${source_zone}" --project "${PROJECT}" "${source_name}":${source_dir}/* "${destination_dir}"
}

function collect-log-mig {
    local group_name="$1"
    local zone="$2"
    local source_dir="$3"
    local destination_dir="$4"
    
    instance_names=()
    instance_names=($(gcloud compute instance-groups managed list-instances \
    "${group_name}" --zone "${zone}" --project "${PROJECT}" \
    --format='value(instance)'))

    for name in "${instance_names[@]}"; do
            collect-log-instance "${name}" "${zone}" "${source_dir}" "${destination_dir}"
    done

}

function copy-logs {
    local vm_name="$1"
    local vm_zone="$2"
    local source_dir="$3"
    local des_dir="$4"
    cmd="mkdir -p ${DES_LOG_DIR}"
    ssh-config "${cmd}" "${vm_name}" "${vm_zone}"
    gcloud compute scp --recurse --zone "${vm_zone}" --project "${PROJECT}" "${DESTINATION}" "${DES_LOG_INSTANCE}":${DES_LOG_DIR}
}


IFS=','; INSTANCE_SERVER_ZONE=($SERVER_ZONE); unset IFS;
IFS=','; INSTANCE_SIM_ZONE=($SIM_ZONE); unset IFS;
IFS=','; INSTANCE_CLIENT_ZONE=($CLIENT_ZONE); unset IFS;
COLLECTDATE="$(date +"%m%d%y-%H%M%S")"
DESTINATION="${DES_LOG_DIR}/${COLLECTDATE}"
if [ ${SERVER_NUM} -gt 0 ]; then
        echo "Collecting logs from ${#INSTANCE_SERVER_ZONE[@]} server machines: "
        if [ ${#INSTANCE_SERVER_ZONE[@]} == 1 ]; then
                collect-log-mig "${SERVER_INSTANCE_PREFIX}-${INSTANCE_SERVER_ZONE[0]}-mig" "${INSTANCE_SERVER_ZONE[0]}" "${SERVER_LOG_DIR}" "${DESTINATION}"
        else
                index=0
                for zone in "${INSTANCE_SERVER_ZONE[@]}"; do
                        collect-log-instance "${SERVER_INSTANCE_PREFIX}-${zone}-${index}" "${zone}" "${SERVER_LOG_DIR}" "${DESTINATION}"
                        index=$(($index + 1)) 
                done

        fi
fi

if [ ${CLIENT_NUM} -gt 0 ]; then
        echo "Collecting logs from ${#INSTANCE_CLIENT_ZONE[@]} client machines: "
        if [ ${#INSTANCE_CLIENT_ZONE[@]} == 1 ]; then
                collect-log-mig "${CLIENT_INSTANCE_PREFIX}-${INSTANCE_CLIENT_ZONE[0]}-mig" "${INSTANCE_CLIENT_ZONE[0]}" "${CLIENT_LOG_DIR}" "${DESTINATION}"
        else
                index=0
                for zone in "${INSTANCE_CLIENT_ZONE[@]}"; do
                        collect-log-instance "${CLIENT_INSTANCE_PREFIX}-${zone}-${index}" "${zone}" "${CLIENT_LOG_DIR}" "${DESTINATION}"
                        index=$(($index + 1)) 
                done

        fi
fi

if [ ${SIM_NUM} -gt 0 ]; then
        echo "Collecting logs from ${#INSTANCE_SIM_ZONE[@]} simulator machines: "
        if [ ${#INSTANCE_SIM_ZONE[@]} == 1 ]; then
                collect-log-mig "${SIM_INSTANCE_PREFIX}-${INSTANCE_SIM_ZONE[0]}-mig" "${INSTANCE_SIM_ZONE[0]}" "${SIM_LOG_DIR}" "${DESTINATION}"
        else
                index=0
                for zone in "${INSTANCE_SIM_ZONE[@]}"; do
                        collect-log-instance "${SIM_INSTANCE_PREFIX}-${zone}-${index}" "${zone}" "${SIM_LOG_DIR}" "${DESTINATION}"
                        index=$(($index + 1)) 
                done

        fi
fi

echo "Copying logs to destination instance."
copy-logs "${DES_LOG_INSTANCE}" "${DES_LOG_INSTANCE_ZONE}" "${DESTINATION}" "${DES_LOG_DIR}"
echo "Removing local copy from ${DESTINATION}"
sudo rm -r "${DESTINATION}"



