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

source scripts/other/set_yum_proxy.sh
source scripts/other/ansible_precheck.sh
source scripts/other/task_log_file.sh
ansible_params=$(grep ^ansible_params: vars/defaults.yml | awk '{ print $2 '} | tr -d '"')

# Run the main playbook
ansible-playbook $ansible_params ./tasks/setup_server.yml -i ./vars/hosts --connection=local
