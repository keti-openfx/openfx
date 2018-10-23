all: fxwatcher fxgateway

fxwatcher:
	protoc -I . \
                -I${GOPATH}/src/github.com/openfx/openfx/pb \
		--go_out=plugins=grpc:. \
		${GOPATH}/src/github.com/openfx/openfx/pb/fxwatcher.proto

fxgateway:
	protoc -I . \
                -I${GOPATH}/src/github.com/openfx/openfx/pb \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:. \
		${GOPATH}/src/github.com/openfx/openfx/pb/fxgateway.proto
	
	protoc -I . \
                -I${GOPATH}/src/github.com/openfx/openfx/pb \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		${GOPATH}/src/github.com/openfx/openfx/pb/fxgateway.proto
	
	protoc -I . \
                -I${GOPATH}/src/github.com/openfx/openfx/pb \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. \
		${GOPATH}/src/github.com/openfx/openfx/pb/fxgateway.proto
	go generate .
