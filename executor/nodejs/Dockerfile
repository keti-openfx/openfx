FROM node:alpine

RUN apk add --no-cache build-base && \
    apk del curl --no-cache

RUN mkdir -p /openfx/watcher
WORKDIR /openfx/watcher

COPY . .

RUN npm install grpc
RUN npm install @grpc/proto-loader

