REGISTRY=10.0.0.180:5000
TAG?=0.1.0

all: proto build push

proto:
	cd pb; $(MAKE)
build:
	docker build -t ${REGISTRY}/fxgateway:$(TAG) .
push:
	docker push ${REGISTRY}/fxgateway:$(TAG)
git:
	git add .
	git commit -m "$m"
	git push origin master
