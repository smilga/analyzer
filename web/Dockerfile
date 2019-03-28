FROM node:10.15.2-stretch-slim

ENV HOST 0.0.0.0
ENV NODE_ENV=production
ENV WEB_DOMAIN={WEB_DOMAIN}

WORKDIR /analyzer-web

COPY . .

RUN npm install --global cross-env
RUN npm install
RUN npm run build
