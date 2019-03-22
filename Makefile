REGISTRY=10.0.0.180:5000
TAG?=0.1.0

all: build push deploy

proto:
	cd pb; $(MAKE)

build:
	go build .

push:
	docker build -t ${REGISTRY}/fxgateway:$(TAG) .
	docker push ${REGISTRY}/fxgateway:$(TAG)

deploy:
#	kubectl delete -f /root/Openfx/yaml/gateway-dep.yml -f /root/Openfx/yaml/gateway-svc.yml
#	kubectl apply -f /root/Openfx/yaml/gateway-dep.yml -f /root/Openfx/yaml/gateway-svc.yml
	kubectl delete -f /root/workspace/openfx/yaml/gateway-dep.yml -f /root/workspace/openfx/yaml/gateway-svc.yml
	kubectl apply -f /root/workspace/openfx/yaml/gateway-dep.yml -f /root/workspace/openfx/yaml/gateway-svc.yml

	#kubectl logs -n openfx $$(kubectl get pods --all-namespaces -l app=fxgateway -o name)
git:
	git add .
	git commit -m "$m"
	git push origin master

clean:
	rm -fr ./openfx-gateway
