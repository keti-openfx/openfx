# Argrumnets for FROM
ARG GO_VERSION=1.9.7

# Build fxwatcher
FROM golang:${GO_VERSION}

RUN mkdir -p ${GOPATH}/src/github.com/keti-openfx/openfx-watcher/go
WORKDIR ${GOPATH}/src/github.com/keti-openfx/openfx-watcher/go

COPY . .

RUN gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*")
RUN go build -o fxwatcher .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
