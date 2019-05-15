# stage: 1
FROM node:latest as react-build

# Create app directory
WORKDIR /app

COPY . ./

RUN yarn install
RUN yarn build

# stage: 2 — the production environment
FROM nginx:alpine
COPY --from=react-build /app/build /usr/share/nginx/html

EXPOSE 80 2222
CMD ["nginx", "-g", "daemon off;"]
