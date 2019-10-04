# Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved.
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

# CCE API
FROM golang:1.12-alpine

ENV GO111MODULE on

WORKDIR /go/src/github.com/otcshare/edgecontroller

RUN apk add git

COPY go.mod .
COPY go.sum .

ARG GITHUB_TOKEN
ENV github_token=$GITHUB_TOKEN
RUN if [[ -n "${github_token}" ]]; \
	then git config --global url.https://${github_token}@github.com/.insteadOf https://github.com/; \
	fi;

RUN go mod download

COPY . .

RUN go build -o dist/cce ./cmd/cce && \
	cp dist/cce /cce

WORKDIR /artifacts

ENTRYPOINT ["/cce"]

CMD []
