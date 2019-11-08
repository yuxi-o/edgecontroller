#!/bin/sh
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

# This script generates key using P-384 curve and certificate for it. Certificate is
# valid for 3 years and signed with CA if CA key and certificate directory is defined

mkdir /root/ca && cp "$3/"* /root/ca/ && cp "$3/cert.pem" "$2/root.pem" && /root/certgen/tls_pair.sh "$1" "$2" /root/ca
