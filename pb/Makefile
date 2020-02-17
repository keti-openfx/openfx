all: fxwatcher fxgateway

fxwatcher:
	protoc -I . \
			-I${GOPATH}/src/github.com/keti-openfx/openfx/pb \
			--go_out=plugins=grpc:. \
			fxwatcher.proto

fxgateway:
	protoc -I . \
				-I${GOPATH}/src/github.com/keti-openfx/openfx/pb \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:. \
		fxgateway.proto
	
	protoc -I . \
				-I${GOPATH}/src/github.com/keti-openfx/openfx/pb \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		fxgateway.proto
	
	protoc -I . \
				-I${GOPATH}/src/github.com/keti-openfx/openfx/pb \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. \
		fxgateway.proto
	go generate .
