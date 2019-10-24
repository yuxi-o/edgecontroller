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

# Add proxy settings to Docker config files

# /root/.docker
folder='/root/.docker'
if ! [[ -d $folder ]]; then
  mkdir $folder
fi

# /root/.docker/config.json
file='/root/.docker/config.json'
if ! [[ -f "$file" ]]; then
  echo '{}' >> $file
fi

for proxy_type in httpProxy httpsProxy noProxy; do
  jq "add(.proxies.default.${proxy_type})" $file | sponge $file
done

# /etc/systemd/system/docker.service.d'
folder='/etc/systemd/system/docker.service.d'
if ! [[ -d $folder ]]; then
  mkdir $folder
fi

# /etc/systemd/system/docker.service.d/http-proxy.conf
file='/etc/systemd/system/docker.service.d/http-proxy.conf'
if ! [[ -f "$file" ]]; then
  echo '[Service]' >> $file
fi

echo "Environment=\"HTTP_PROXY=${1}/\"" >> $file 
echo "Environment=\"HTTPS_PROXY=${2}/\"" >> $file 
echo "Environment=\"NO_PROXY=${3}/\"" >> $file 
