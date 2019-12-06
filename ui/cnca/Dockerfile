# SPDX-License-Identifier: Apache-2.0
# Copyright © 2019 Intel Corporation

# stage: 1
FROM node:lts-alpine as react-build

# Set the AF API URL variable during build time
ARG REACT_APP_CNCA_AF_API
ENV REACT_APP_CNCA_AF_API=$REACT_APP_CNCA_AF_API
ARG REACT_APP_CNCA_5GOAM_API
ENV REACT_APP_CNCA_5GOAM_API=$REACT_APP_CNCA_5GOAM_API

# Set the Controller UI URL variable during build time
ARG CONTROLLER_UI_URL
ENV REACT_APP_CONTROLLER_UI_URL=$CONTROLLER_UI_URL

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
#
RUN yarn install --production

EXPOSE 80
CMD ["node", "server.js"]
