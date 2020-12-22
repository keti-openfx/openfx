REGISTRY=10.0.0.91:5000
TAG?=0.1.0

all: build

build:
	echo "#### Compiling Openfx OAuth 2.0 Server  ####"
	docker build  --network=host -t ${REGISTRY}/oauth2-server:$(TAG) .
	docker push ${REGISTRY}/oauth2-server:$(TAG)

clean:
	@rm -fr ./oauth2

.PHONY: all build clean
