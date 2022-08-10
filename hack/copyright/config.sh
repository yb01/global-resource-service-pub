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
GRS_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..


COPYRIGHT_MATCH=${COPYRIGHT_MATCH:-"limitations under the License"}
COPYRIGHT_YEAR=${COPYRIGHT_YEAR:-"$( date +%Y )"}
GRS_COPYRIGHT_LINE_MODIFIED_GO=${GRS_COPYRIGHT_LINE_MODIFIED_GO:-"Copyright ${COPYRIGHT_YEAR} Authors of Global Resource Service - file modified."}
GRS_COPYRIGHT_LINE_MODIFIED_OTHER=${GRS_COPYRIGHT_LINE_MODIFIED_OTHER:-"# Copyright ${COPYRIGHT_YEAR} Authors of Global Resource Service - file modified."}
K8S_COPYRIGHT_MATCH=${K8S_COPYRIGHT_MATCH:-"The Kubernetes Authors"}
GRS_COPYRIGHT_MATCH=${GRS_COPYRIGHT_MATCH:-"Authors of Global Resource Service"}
GRS_COPYRIGHT_LINE_NEW_GO="/*\nCopyright ${COPYRIGHT_YEAR} Authors of Global Resource Service.\n\nLicensed under the Apache License, Version 2.0 (the \"License\");\nyou may not use this file except in compliance with the License.\nYou may obtain a copy of the License at\n\n    http://www.apache.org/licenses/LICENSE-2.0\n\nUnless required by applicable law or agreed to in writing, software\ndistributed under the License is distributed on an \"AS IS\" BASIS,\nWITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\nSee the License for the specific language governing permissions and\nlimitations under the License.\n*/\n"
GRS_COPYRIGHT_LINE_NEW_OTHER="#\n# Copyright ${COPYRIGHT_YEAR} Authors of Global Resource Service.\n#\n# Licensed under the Apache License, Version 2.0 (the \"License\");\n# you may not use this file except in compliance with the License.\n# You may obtain a copy of the License at\n#\n#     http://www.apache.org/licenses/LICENSE-2.0\n#\n# Unless required by applicable law or agreed to in writing, software\n# distributed under the License is distributed on an \"AS IS\" BASIS,\n# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n# See the License for the specific language governing permissions and\n# limitations under the License.\n"


DAY0_COMMIT=$( git rev-list --max-parents=0 HEAD | tail -n 1 )
MERGED_COMMIT=$( git log --show-signature --oneline | grep "gpg: Signature made" | head -n 1 | cut -c1-7 )
START_COMMIT=${START_COMMIT:-"${MERGED_COMMIT}"}
GRS_REPO=${GRS_REPO:-"https://github.com/CentaurusInfra/global-resource-service"}
