FROM node:lts-slim AS build
WORKDIR /src
RUN npm install -g @angular/cli

COPY package*.json ./
RUN npm ci

COPY . ./

RUN echo "environment contents:" && cat src/environments/environment.prod.ts

RUN ng build -c production

FROM nginx:stable AS final
EXPOSE 4200
COPY --from=build src/dist/angular-app/browser  /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf