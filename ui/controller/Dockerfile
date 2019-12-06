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

# stage: 1
FROM node:lts-alpine as react-build

# Set the Controller API URL variable during build time
ARG REACT_APP_CONTROLLER_API
ENV REACT_APP_CONTROLLER_API=$REACT_APP_CONTROLLER_API

# Set the Controller orchastration mode variable during build time
ARG CCE_ORCHESTRATION_MODE
ENV REACT_APP_ORCHESTRATION_MODE=$CCE_ORCHESTRATION_MODE

# Set the UI URL variables during build time
ARG CONTROLLER_UI_URL
ENV REACT_APP_CONTROLLER_UI_URL=$CONTROLLER_UI_URL

ARG CNCA_UI_URL
ENV REACT_APP_CNCA_UI_URL=$CNCA_UI_URL

ARG CUPS_UI_URL
ENV REACT_APP_CUPS_UI_URL=$CUPS_UI_URL

# Instruct webpack to only import runtime scripts
ENV INLINE_RUNTIME_CHUNK=false

# Create app directory
WORKDIR /app

# Copy source
COPY . ./

# Install production dependencies for app
RUN yarn install --production

# Build production bundle
RUN yarn build

# stage: 2 — the production environment
FROM node:lts-alpine

# Create app directory
WORKDIR /usr/src/app

# Copy production bundle from build stage
COPY --from=react-build /app/build ./build

# Copy static server
COPY ./server ./

# Install production dependencies for server
RUN yarn install --production

EXPOSE 80
CMD ["node", "server.js"]
