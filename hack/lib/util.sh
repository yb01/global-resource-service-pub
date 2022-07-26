#!/usr/bin/env bash

function ssh-config {
        local cmd="$1"
        local machine_name="$2"
        local zone="$3"
        gcloud compute ssh \
        "${machine_name}" \
        --ssh-flag="-o LogLevel=quiet" \
        --ssh-flag="-o ConnectTimeout=30" \
        --project "${PROJECT}" \
        --zone "${zone}" \
        --command "${cmd}" \
        --quiet &
}