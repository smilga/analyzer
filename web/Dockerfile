FROM node:10.15.2-stretch-slim

ENV HOST 0.0.0.0
ENV NODE_ENV=production

WORKDIR /analyzer-web

COPY . .

RUN npm install --global cross-env
RUN npm install
