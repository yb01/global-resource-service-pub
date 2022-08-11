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

if [ "${DESTINATION}" == "" ]; then
  echo "Env: DESTINATION cannot be empty, Please double check."
  exit 1
fi
mkdir -p ${DESTINATION}/csv

csv_name="result.csv"
if [ "${COLLECTDATE}" != "" ]; then
  csv_name="result-${COLLECTDATE}.csv"
fi


function grep-string {
  local file_name="$1"
  local start_string="$2"
  local end_string="${3:-}"

  grep_string=$(grep "${start_string}" ${file_name})
  if [ "${end_string}" == "" ]; then
    echo "$grep_string" | sed -E "s/.*${start_string}//"
  else
    echo "$grep_string" | sed -E "s/.*${start_string}(.*)${end_string}.*/\1/"
  fi

}

cd ${DESTINATION}

echo "Collecting scheduler test summary to csv"
echo "file name","RegisterClientDuration","ListDuration","Number of nodes listed","Watch session last","Number of nodes Added","Updated","Deleted","watch prolonged than 1s","Watch perc50","Watch perc90","Watch perc99">> ./csv/${csv_name}
for name in $( ls | grep client);do
  start_string="RegisterClientDuration: "
  end_string=""
  register_client_duration=$(grep-string "${name}" "${start_string}" "${end_string}")
  
  start_string="ListDuration: "
  end_string=". Number"
  list_duration=$(grep-string "${name}" "${start_string}" "${end_string}")
  
  start_string="Number of nodes listed: "
  end_string=""
  nodes_listed=$(grep-string "${name}" "${start_string}" "${end_string}")
  
  start_string="Watch session last: "
  end_string=". Number"
  watch_session_last=$(grep-string "${name}" "${start_string}" "${end_string}")

  start_string="Number of nodes Added :"
  end_string=", Updated"
  number_nodes_added=$(grep-string "${name}" "${start_string}" "${end_string}") 
  
  start_string="Updated: "
  end_string=", Deleted"
  number_nodes_updated=$(grep-string "${name}" "${start_string}" "${end_string}")
  
  start_string="Deleted: "
  end_string=". watch prolonged"
  number_nodes_deleted=$(grep-string "${name}" "${start_string}" "${end_string}") 
  
  start_string="watch prolonged than 1s: "
  end_string=""
  watch_prolonged_than1s=$(grep-string "${name}" "${start_string}" "${end_string}") 
  
  start_string="perc50 "
  end_string=", perc90"
  watch_perc50=$(grep-string "${name}" "${start_string}" "${end_string}")
  
  start_string="perc90 "
  end_string=", perc99"
  watch_perc90=$(grep-string "${name}" "${start_string}" "${end_string}")
  
  start_string="perc99 "
  end_string=". Total"
  watch_perc99=$(grep-string "${name}" "${start_string}" "${end_string}")

  echo "${name}","${register_client_duration}","${list_duration}","${nodes_listed}","${watch_session_last}","${number_nodes_added}","${number_nodes_updated}","${number_nodes_deleted}","${watch_prolonged_than1s}","${watch_perc50}","${watch_perc90}","${watch_perc99}" >> ./csv/${csv_name}
done

###adding empty line to csv
echo "" >> ./csv/${csv_name}
echo "" >> ./csv/${csv_name}
echo "" >> ./csv/${csv_name}
echo "" >> ./csv/${csv_name}


echo "Collecting service test summary to csv"
for name in $( ls | grep server);do
  echo "${name}">> ./csv/${csv_name}
  echo "Metrics Item","perc50","perc90","perc99","Total count">> ./csv/${csv_name}
  metrics_item="[Metrics][AGG_RECEIVED]"
  grep_line=$(grep "\[Metrics\]\[AGG_RECEIVED\]" $name | tail -1)
  perc50=$(echo $grep_line | sed "s/.*perc50 //; s/, perc90.*//")
  perc90=$(echo $grep_line | sed "s/.*perc90 //; s/, perc99.*//")
  perc99=$(echo $grep_line | sed "s/.*perc99 //; s/\. Total.*//")
  total_count=$(echo $grep_line | sed "s/.*Total count //")
  echo "${metrics_item}","${perc50}","${perc90}","${perc99}","${total_count}">> ./csv/${csv_name}
  
  metrics_item="[Metrics][DIS_RECEIVED]"
  grep_line=$(grep "\[Metrics\]\[DIS_RECEIVED\]" $name | tail -1)
  perc50=$(echo $grep_line | sed "s/.*perc50 //; s/, perc90.*//")
  perc90=$(echo $grep_line | sed "s/.*perc90 //; s/, perc99.*//")
  perc99=$(echo $grep_line | sed "s/.*perc99 //; s/\. Total.*//")
  total_count=$(echo $grep_line | sed "s/.*Total count //")
  echo "${metrics_item}","${perc50}","${perc90}","${perc99}","${total_count}">> ./csv/${csv_name}
  
  metrics_item="[Metrics][DIS_SENDING]"
  grep_line=$(grep "\[Metrics\]\[DIS_SENDING\]" $name | tail -1)
  perc50=$(echo $grep_line | sed "s/.*perc50 //; s/, perc90.*//")
  perc90=$(echo $grep_line | sed "s/.*perc90 //; s/, perc99.*//")
  perc99=$(echo $grep_line | sed "s/.*perc99 //; s/\. Total.*//")
  total_count=$(echo $grep_line | sed "s/.*Total count //")
  echo "${metrics_item}","${perc50}","${perc90}","${perc99}","${total_count}">> ./csv/${csv_name}

  metrics_item="[Metrics][DIS_SENT]"
  grep_line=$(grep "\[Metrics\]\[DIS_SENT\]" $name | tail -1)
  perc50=$(echo $grep_line | sed "s/.*perc50 //; s/, perc90.*//")
  perc90=$(echo $grep_line | sed "s/.*perc90 //; s/, perc99.*//")
  perc99=$(echo $grep_line | sed "s/.*perc99 //; s/\. Total.*//")
  total_count=$(echo $grep_line | sed "s/.*Total count //")
  echo "${metrics_item}","${perc50}","${perc90}","${perc99}","${total_count}">> ./csv/${csv_name}

  metrics_item="[Metrics][SER_ENCODED]"
  grep_line=$(grep "\[Metrics\]\[SER_ENCODED\]" $name | tail -1)
  perc50=$(echo $grep_line | sed "s/.*perc50 //; s/, perc90.*//")
  perc90=$(echo $grep_line | sed "s/.*perc90 //; s/, perc99.*//")
  perc99=$(echo $grep_line | sed "s/.*perc99 //; s/\. Total.*//")
  total_count=$(echo $grep_line | sed "s/.*Total count //")
  echo "${metrics_item}","${perc50}","${perc90}","${perc99}","${total_count}">> ./csv/${csv_name}

  metrics_item="[Metrics][SER_SENT]"
  grep_line=$(grep "\[Metrics\]\[SER_SENT\]" $name | tail -1)
  perc50=$(echo $grep_line | sed "s/.*perc50 //; s/, perc90.*//")
  perc90=$(echo $grep_line | sed "s/.*perc90 //; s/, perc99.*//")
  perc99=$(echo $grep_line | sed "s/.*perc99 //; s/\. Total.*//")
  total_count=$(echo $grep_line | sed "s/.*Total count //")
  echo "${metrics_item}","${perc50}","${perc90}","${perc99}","${total_count}">> ./csv/${csv_name}

done

echo "Please check generated csv report under ${DESTINATION}/csv/${csv_name}"