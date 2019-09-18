registry=10.0.0.180:5000
version=0.1.0
go_version=1.9.7

all: proto docker-build docker-push

proto:
	protoc -I . \
                -I${GOPATH}/src/github.com/keti-openfx/openfx/pb \
                --go_out=plugins=grpc:./pb \
                ${GOPATH}/src/github.com/keti-openfx/openfx/pb/fxwatcher.proto
build:
	go build -o fxwatcher .
docker-build:
	docker build -t $(registry)/fxwatcher:$(version)-go --build-arg GO_VERSION=$(go_version) .
docker-push:
	docker push $(registry)/fxwatcher:$(version)-go
