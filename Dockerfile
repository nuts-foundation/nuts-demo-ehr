FROM node:13-alpine

LABEL maintainer="steven.vandervegt@nuts.nl"

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY ./ .

EXPOSE 8000 8001 8002

ENTRYPOINT npm start

