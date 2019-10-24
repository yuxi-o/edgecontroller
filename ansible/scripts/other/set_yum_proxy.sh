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

set -e

echo "Proxy check..."

# Remove yum proxy for initial checkings
sed "/^proxy=.*/d" -i /etc/yum.conf

# Set proxy for yum if defines and proxy enabled in Ansible ./vars/defaults.yml
proxy_enabled=$(grep ^proxy_enable: vars/defaults.yml | tr -d "\ |'|\"" | cut -d: -f2-)
proxy_yum=$(grep ^proxy_yum: vars/defaults.yml | tr -d "\ |'|\"" | cut -d: -f2-)
if [[ "$proxy_enabled" == *"true"* ]] && [[ "$proxy_yum" != "" ]]; then
  echo "Using proxy for yum: ${proxy_yum}/"
  echo "proxy=${proxy_yum}/" >> /etc/yum.conf
else
  echo "Not using proxy"
fi

echo "done"
