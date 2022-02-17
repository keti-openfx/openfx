FROM openjdk:8

RUN apt-get update 

RUN mkdir -p /openfx/executor
WORKDIR /openfx/executor

COPY . .
RUN chmod +x gradlew

