# Build fxwatcher
FROM golang:1.16

ENV GO111MODULE=on

RUN mkdir -p ${GOPATH}/src/github.com/keti-openfx/openfx/executor/go
WORKDIR ${GOPATH}/src/github.com/keti-openfx/openfx/executor/go

COPY . .

RUN go get github.com/minio/minio-go
RUN go build -o fxwatcher .

