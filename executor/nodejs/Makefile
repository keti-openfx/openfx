registry=keti.asuscomm.com:5000
version=0.1.0
node_version=13.1.0

all: docker-build docker-push

docker-build:
	docker build --network=host -t $(registry)/fxwatcher:$(version)-nodejs --build-arg NODE_VERSION=$(node_version) .

docker-push:
	docker push $(registry)/fxwatcher:$(version)-nodejs

