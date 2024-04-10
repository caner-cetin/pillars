FROM node:latest as develop-stage
WORKDIR /frontend
RUN yarn config set ignore-engines true
COPY ./frontend/package.json ./
COPY ./frontend/quasar.config.js ./
RUN yarn global add @quasar/cli
COPY ./frontend/ ./

# build stage
FROM develop-stage as build-stage
RUN yarn install
RUN quasar build

# production stage
FROM nginx:latest as production-stage
COPY --from=build-stage /frontend/dist/spa /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]