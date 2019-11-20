FROM node:alpine

RUN apk add --no-cache build-base && \
    apk del curl --no-cache

RUN npm install grpc
RUN npm install @grpc/proto-loader

RUN mkdir -p /openfx/watcher
WORKDIR /openfx/watcher

COPY . .

