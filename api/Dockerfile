FROM golang:1.11.0-alpine3.8

WORKDIR /analyzer-api

COPY . .

RUN apk update && apk upgrade && \
        apk add nodejs-current && \
        apk add nodejs-npm && \
        apk add --no-cache git

WORKDIR /analyzer-api/puppeteer

RUN npm install

WORKDIR /analyzer-api

RUN CGO_ENABLED=0 go install -v ./...
