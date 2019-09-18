registry=10.0.0.180:5000
version=0.1.0

all: proto python2 python3

clean:
	rm -f ./*pb2.py ./*pb2_grpc.py

proto:
	python2 -m grpc_tools.protoc \
		-I . \
		-I${GOPATH}/src/github.com/keti-openfx/openfx/pb \
		--python_out=. \
		--grpc_python_out=. \
		${GOPATH}/src/github.com/keti-openfx/openfx/pb/fxwatcher.proto


python2: python_version=2.7
python2:
	docker build -t $(registry)/fxwatcher:$(version)-python2 --build-arg PYTHON_VERSION=$(python_version) -f Dockerfile-py2 .
	docker push $(registry)/fxwatcher:$(version)-python2

python3: python_version=3.4
python3:
	docker build -t $(registry)/fxwatcher:$(version)-python3 --build-arg PYTHON_VERSION=$(python_version) -f Dockerfile-py3 .
	docker push $(registry)/fxwatcher:$(version)-python3
