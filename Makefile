REGISTRY=10.0.0.180:5000
TAG?=0.1.0

all: build push

proto:
	cd pb; $(MAKE)

build:
	go build .

push:
	docker build -t ${REGISTRY}/fxgateway:$(TAG) .
	docker push ${REGISTRY}/fxgateway:$(TAG)

deploy:
	kubectl apply -f /root/Openfx/yaml/gateway-dep.yml
	kubectl apply -f /root/Openfx/yaml/gateway-svc.yml
git:
	git add .
	git commit -m "$m"
	git push origin master

clean:
	rm -fr ./openfx-gateway
