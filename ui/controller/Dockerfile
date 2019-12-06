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

# Build image with online dependencies
FROM node:lts-alpine AS cce-ui-deps-image

# Create app directory & install dependencies for running the server
WORKDIR /usr/src/app
COPY ./server/package.json ./server/yarn.lock ./
RUN yarn install --pure-lockfile --production
COPY ./server ./

# Create app directory & install dependencies for building the app
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --pure-lockfile --production

# Second stage - building the app with already existing dependencies
FROM cce-ui-deps-image

# Set the Controller envs for UI build
ARG CCE_ORCHESTRATION_MODE
ENV REACT_APP_ORCHESTRATION_MODE=$CCE_ORCHESTRATION_MODE

ARG REACT_APP_CONTROLLER_API
ENV REACT_APP_CONTROLLER_API=$REACT_APP_CONTROLLER_API

ARG CONTROLLER_UI_URL
ENV REACT_APP_CONTROLLER_UI_URL=$CONTROLLER_UI_URL

ARG CNCA_UI_URL
ENV REACT_APP_CNCA_UI_URL=$CNCA_UI_URL

ARG CUPS_UI_URL
ENV REACT_APP_CUPS_UI_URL=$CUPS_UI_URL

# Instruct webpack to only import runtime scripts
ENV INLINE_RUNTIME_CHUNK=false

# Copy UI source & build
COPY . ./
RUN yarn build

WORKDIR /usr/src/app
RUN cp -r /app/build ./build

EXPOSE 80
CMD ["node", "server.js"]
