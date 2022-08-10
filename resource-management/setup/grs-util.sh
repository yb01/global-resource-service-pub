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


GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

CLOUD_PROVIDER="${CLOUD_PROVIDER:-gce}"

# PROVIDER_VARS is a list of cloud provider specific variables. Note:
# this is a list of the _names_ of the variables, not the value of the
# variables.

PROVIDER_UTILS="${GRS_ROOT}/setup/${CLOUD_PROVIDER}/util.sh"
if [ -f "${PROVIDER_UTILS}" ]; then
    source "${PROVIDER_UTILS}"
fi
