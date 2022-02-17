registry=10.0.0.255:5000/openfx
version=1.0
csharp_version=2.1

all: docker-build docker-push

docker-build:
	docker build --network=host -t $(registry)/fxwatcher:$(version)-csharp --build-arg CSHARP_VERSION=$(csharp_version) .

docker-push:
	docker push $(registry)/fxwatcher:$(version)-csharp
