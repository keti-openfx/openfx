###############################################
#### Edit to your private registry address ####
###############################################
registry=keti.dev:5000
version=0.1.0

all: docker-build docker-push

docker-build:
	docker build --network=host -t $(registry)/fxwatcher:$(version)-nodejs .

docker-push:
	docker push $(registry)/fxwatcher:$(version)-nodejs

