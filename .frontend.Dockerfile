FROM node:latest as develop-stage
WORKDIR /frontend
RUN yarn config set ignore-engines true
COPY ./frontend/package.json ./
COPY ./frontend/quasar.config.js ./
RUN yarn install --production=false
COPY ./frontend/ ./
RUN yarn global add @quasar/cli --production=false
EXPOSE 8080
RUN quasar update -i
ENTRYPOINT ["quasar", "dev", "--p", "8080"]
# FROM nginx:latest as production-stage
# COPY --from=build-stage /frontend/dist/spa /usr/share/nginx/html
# EXPOSE 80
# CMD ["nginx", "-g", "daemon off;"]