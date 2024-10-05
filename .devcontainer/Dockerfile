FROM golang:alpine

RUN apk update && apk add --no-cache \
  make \
  git \
  curl \
  gcc \
  make \ 
  build-base

RUN rm -rf /var/cache/apk/*