# Build fxwatcher
FROM golang:1.12

ENV GO111MODULE=on

RUN mkdir -p ${GOPATH}/src/github.com/keti-openfx/openfx/executor/go
WORKDIR ${GOPATH}/src/github.com/keti-openfx/openfx/executor/go

COPY . .

RUN go mod vendor -v

RUN gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*")
RUN go build -o fxwatcher .

