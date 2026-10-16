FROM node:18-alpine AS build-env

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

ARG VUE_APP_VERSION
ENV VUE_APP_VERSION=$APP_VERSION

RUN npm run build


FROM nginx:alpine AS release

COPY --from=build-env /app/dist/ /usr/share/nginx/html/
COPY ./nginx.conf /etc/nginx/conf.d/default.conf