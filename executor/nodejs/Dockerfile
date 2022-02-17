FROM node:alpine

RUN apk add --no-cache build-base && \
    apk del curl --no-cache

RUN mkdir -p /openfx/executor
WORKDIR /openfx/executor

COPY . .

RUN npm install @grpc/grpc-js
RUN npm install @grpc/proto-loader

