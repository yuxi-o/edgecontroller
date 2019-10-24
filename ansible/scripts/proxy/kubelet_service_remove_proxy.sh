#!/usr/bin/env bash
#
# Copyright 2019 Intel Corporation. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Remove proxy from kubelet system service file

# /usr/lib/systemd/system/kubelet.service
file='/usr/lib/systemd/system/kubelet.service'
if [[ -f "$file" ]]; then
  for proxy_var in HTTP_PROXY HTTPS_PROXY NO_PROXY; do
    sed "/^Environment=\"${proxy_var}.*/d" -i $file
  done
fi
