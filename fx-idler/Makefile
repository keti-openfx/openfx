REGISTRY=10.0.0.52:5000
TAG?=0.1.0


all: build push 

build:
	@echo "#### Compiling OpenFx Fx-Idler####"
	@go build .

push:
	docker build --network=host -t ${REGISTRY}/fx-idler:$(TAG) .
	docker push ${REGISTRY}/fx-idler:$(TAG)