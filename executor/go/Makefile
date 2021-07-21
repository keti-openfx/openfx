registry=10.0.0.202:5000
version=2.0
go_version=1.12

all: docker-build docker-push

proto:
	protoc -I . \
                -I${GOPATH}/src/github.com/keti-openfx/openfx/pb \
                --go_out=plugins=grpc:./pb \
                ${GOPATH}/src/github.com/keti-openfx/openfx/pb/fxwatcher.proto
build:
	go build -o fxwatcher .

docker-build:
	docker build --network=host -t $(registry)/fxwatcher:$(version)-go --build-arg GO_VERSION=$(go_version) .

docker-push:
	docker push $(registry)/fxwatcher:$(version)-go
