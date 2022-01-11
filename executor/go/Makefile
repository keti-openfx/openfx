###############################################
#### Edit to your private registry address ####
###############################################
registry=keti.dev:5000
version=0.1.0

all: docker-build docker-push

proto:
	protoc -I . \
                -I${GOPATH}/src/github.com/keti-openfx/openfx/pb \
                --go_out=plugins=grpc:./pb \
                ${GOPATH}/src/github.com/keti-openfx/openfx/pb/fxwatcher.proto
build:
	go build -o fxwatcher .

docker-build:
	docker build --network=host -t $(registry)/fxwatcher:$(version)-go .

docker-push:
	docker push $(registry)/fxwatcher:$(version)-go
