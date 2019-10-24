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

# Remove proxy settings from OS and Docker config files

# /etc/environment
for proxy_line in http_proxy HTTP_PROXY https_proxy HTTPS_PROXY ftp_proxy FTP_PROXY no_proxy NO_PROXY; do
  sed "/^${proxy_line}.*/d" -i /etc/environment
done

# /etc/yum.conf
file='/etc/yum.conf'
if [[ -f "$file" ]]; then
  sed "/^proxy=.*/d" -i $file
fi
