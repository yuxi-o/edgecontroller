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

COPY conf/default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
