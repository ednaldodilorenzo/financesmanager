FROM node:18-slim AS development

WORKDIR /app

# Install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Start development server
CMD ["npm", "run", "serve"]


FROM node:18-alpine AS build-env

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build


FROM nginx:alpine as release

COPY --from=build-env /app/dist/ /usr/share/nginx/html/
COPY ./nginx.conf /etc/nginx/conf.d/default.conf