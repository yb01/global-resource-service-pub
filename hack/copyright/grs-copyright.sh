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

GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
source "${GRS_ROOT}/hack/copyright/config.sh"

TMPDIR="/tmp/grsCopyright"
REPODIRNAME="${REPODIRNAME:-$GRS_ROOT}"
LOGFILENAME="grsCopyrightTool.log"
LOGDIR=$TMPDIR
LOGFILE=$LOGDIR/$LOGFILENAME
CHANGEDFILES="$LOGDIR/changed_files"
EXIT_ERROR=0

SED_CMD=""
STAT_CMD=""
TOUCH_CMD=""

if [ -d "${LOGDIR}" ]; then
    sudo rm ${LOGDIR}/*
else
    mkdir -p ${LOGDIR}
fi

if [[ "$OSTYPE" == "darwin"* ]]
then
    SED_CMD=`which gsed`
    if [ -z $SED_CMD ]
    then
        echo "Please install gnu-sed (brew install gnu-sed)"
        exit 1
    fi
    STAT_CMD="stat -f %Sm -t %Y%m%d%H%M.%S "
    TOUCH_CMD="touch -mt "
elif [[ "$OSTYPE" == "linux"* ]]
then
    SED_CMD=`which sed`
    if [ -z $SED_CMD ]
    then
        echo "Please install sed"
        exit 1
    fi
    STAT_CMD="stat -c %y "
    TOUCH_CMD="touch -d "
else
    echo "Unsupported OS $OSTYPE"
    exit 1
fi

display_usage() {
    echo "Usage: $0 <optional-grs-repo-path> <optional-log-directory>"
    echo "       If optional global-resource-service repo path is provided, repo setup step will be skipped"
}

if [ ! -z $2 ]
then
    LOGDIR=$2
    LOGFILE=$LOGDIR/$LOGFILENAME
fi

if [ ! -z $1 ]
then
    if [[ ( $1 == "--help") ||  $1 == "-h" ]]
    then
        display_usage
        exit 0
    else
        REPODIRNAME=$1
        if [ -z $2 ]
        then
	    LOGFILE=$REPODIRNAME/../$LOGFILENAME
        fi
        rm -f $LOGFILE
        inContainer=true
        if [[ -f /proc/1/sched ]]
        then
            PROC1=`cat /proc/1/sched | head -n 1`
            if [[ $PROC1 == systemd* ]]
            then
                inContainer=false
            fi
        else
            if [[ "$OSTYPE" == "darwin"* ]]
            then
                inContainer=false
            fi
        fi
        if [ "$inContainer" = true ]
        then
            echo "WARN: Skipping copyright check for in-container build as git repo is not available"
            echo "WARN: Skipping copyright check for in-container build as git repo is not available" >> $LOGFILE
            exit 0
        else
            echo "Running copyright check for repo: $REPODIRNAME, logging to $LOGFILE"
        fi
    fi
fi

function clone_repo() {
    local REPO=$1
    local DESTDIR=$2
    git clone $REPO $DESTDIR
}

function setup_repos() {
    if [ -d $TMPDIR ]; then
        rm -rf $TMPDIR
    fi
    mkdir -p $TMPDIR
    clone_repo $GRS_REPO $REPODIRNAME
}


function update-copyright() {
    local repo_file=$1
    set +e
    cat $repo_file | grep "$COPYRIGHT_MATCH" > /dev/null 2>&1
#    local tstamp=$($STAT_CMD $repo_file)
    if [ $? -eq 0 ]; then 
        cat $repo_file | grep "$GRS_COPYRIGHT_MATCH" > /dev/null 2>&1
        if [ $? -eq 0 ]; then
            echo "WARN: File $repo_file already has grs copyright. Ignoring." >> $LOGFILE
        else
            cat $repo_file | grep "$K8S_COPYRIGHT_MATCH" > /dev/null 2>&1
            if [ $? -eq 0 ]; then
                echo "WARN: File $repo_file has k8s copyright. Adding grs copyright." >> $LOGFILE
                if [[ $repo_file = *.go ]] || [[ $repo_file = *.proto ]]; then
                    $SED_CMD -i "/${K8S_COPYRIGHT_MATCH}/a ${GRS_COPYRIGHT_LINE_MODIFIED_GO}" $repo_file
                else
                    $SED_CMD -i "/${K8S_COPYRIGHT_MATCH}/a ${GRS_COPYRIGHT_LINE_MODIFIED_OTHER}" $repo_file
                fi
            else
                echo "WARN: File $repo_file does not have either K8s or grs copyright. Adding grs copyright." >> $LOGFILE
                if [[ $repo_file = *.go ]] || [[ $repo_file = *.proto ]]; then
                    $SED_CMD -i "1 i ${GRS_COPYRIGHT_LINE_NEW_GO}" $repo_file
                else
                    if [[ $repo_file = *.sh ]]; then
                        bash_after="\#\!\/usr\/bin"
                        cat $repo_file | grep "$bash_after" > /dev/null 2>&1
                        if [ $? -eq 0 ]; then
                            $SED_CMD -i "/${bash_after}/a ${GRS_COPYRIGHT_LINE_NEW_OTHER}" $repo_file
                        else
                            $SED_CMD -i "1 i  ${GRS_COPYRIGHT_LINE_NEW_OTHER}" $repo_file
                        fi
                    fi
                fi
            fi
        fi
    else  
        echo "Warning: File $repo_file does not have either K8s or grs copyright. Adding grs copyright." >> $LOGFILE
        if [[ $repo_file = *.go ]] || [[ $repo_file = *.proto ]]; then
            $SED_CMD -i "1 i ${GRS_COPYRIGHT_LINE_NEW_GO}" $repo_file
        else
            if [[ $repo_file = *.sh ]]; then
                bash_after="\#\!\/usr\/bin"
                cat $repo_file | grep "$bash_after" > /dev/null 2>&1
                if [ $? -eq 0 ]; then
                    $SED_CMD -i "/${bash_after}/a ${GRS_COPYRIGHT_LINE_NEW_OTHER}" $repo_file
                else
                    $SED_CMD -i "1 i  ${GRS_COPYRIGHT_LINE_NEW_OTHER}" $repo_file
                fi
            fi
        fi
    fi
#    $TOUCH_CMD "$tstamp" $repo_file
    set -e
}

function update-grs-copyrights() {
    echo "Inspecting changed files and update grs copyright, writing logs to $LOGFILE"
    local changed_files=$1
    while IFS= read -r line
    do
        if [ ! -z $line ]
        then
            update-copyright $REPODIRNAME/$line
        fi
    done < $changed_files
    echo "Done."
}

function get_changed_files_list() {
    pushd $REPODIRNAME
    local start_commit=$1
    local head_commit=`git rev-list HEAD | head -n 1`
    echo "Getting changed files between start_commit: $start_commit, head_commit: $head_commit"
    git diff --name-only $start_commit $head_commit | \
        egrep -v "\.git|\.md|\.bazelrc|\.json|\.pb|\.yaml|BUILD|boilerplate|vendor\/" | \
        egrep -v "\.mod|\.sum|\.png|\.PNG|OWNERS" >> $CHANGEDFILES || true
    popd
    echo "Saved changed_files at $CHANGEDFILES"
}

get_changed_files_list ${START_COMMIT}
update-grs-copyrights $CHANGEDFILES

exit $EXIT_ERROR
