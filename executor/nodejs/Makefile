registry=keti.asuscomm.com:5000
version=0.1.0

all: docker-build docker-push

docker-build:
	docker build -t $(registry)/fxwatcher:$(version)-nodejs .

docker-push:
	docker push $(registry)/fxwatcher:$(version)-nodejs

