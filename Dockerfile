FROM node:13-alpine

LABEL maintainer="steven.vandervegt@nuts.nl"

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY ./ .

EXPOSE 80 81 82

ENTRYPOINT npm start

