# stage: 1
FROM node:lts-alpine as react-build

# Set the Controller API URL variable during build time
ARG REACT_APP_CUPS_API
ENV REACT_APP_CUPS_API=$REACT_APP_CUPS_API

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
