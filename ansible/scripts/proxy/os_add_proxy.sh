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

# Remove proxy settings from OS config files

# /etc/environment
file='/etc/environment'
echo "http_proxy=${proxy_http}" >> $file
echo "HTTP_proxy=${proxy_http}" >> $file
echo "https_proxy=${proxy_https}" >> $file
echo "HTTPS_proxy=${proxy_https}" >> $file
echo "ftp_proxy=${proxy_ftp}" >> $file
echo "FTP_proxy=${proxy_ftp}" >> $file
echo "no_proxy=${proxy_noproxy}" >> $file
echo "NO_PROXY=${proxy_noproxy}" >> $file

# /etc/yum.conf
file='/etc/yum.conf'
echo "proxy=${proxy_yum}" >> $file
