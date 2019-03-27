FROM golang:1.12.1-alpine3.9

WORKDIR /analyzer-api

COPY . .

RUN apk update && apk upgrade && apk add --no-cache git

RUN CGO_ENABLED=0 go install -v ./...
