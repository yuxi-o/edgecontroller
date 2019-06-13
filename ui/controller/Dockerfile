# stage: 1
FROM node:lts-alpine as react-build

# Create app directory
WORKDIR /app

# Copy source
COPY . ./

# Install production dependencies for app
RUN yarn install --production

# Instruct webpack to only import runtime scripts
ENV INLINE_RUNTIME_CHUNK=false

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

# Copy environment files
COPY .env.production ./.env

# Install production dependencies for server
RUN yarn install --production

EXPOSE 80
CMD ["node", "server.js"]
