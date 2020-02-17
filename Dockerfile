FROM golang:1.12 as builder

RUN mkdir -p /go/src/github.com/keti-openfx/openfx

WORKDIR /go/src/github.com/keti-openfx/openfx

COPY . .

ENV GO111MODULE=on
RUN gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*")

#WORKDIR /go/src/github.com
#RUN go get -u google.golang.org/grpc@v1.13.0
#RUN go get -u golang.org/x/sys/unix
#RUN go get github.com/golang/protobuf@v1.1.0
#RUN git clone https://github.com/grpc-ecosystem/grpc-gateway.git

#WORKDIR /go/src/github.com/grpc-ecosystem/grpc-gateway
#RUN make
#RUN make install
#RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.4.1
#RUN go get -u github.com/golang/protobuf/protoc-gen-go
#RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.4.1

#RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
#RUN apt-get update && apt-get install -y unzip
#RUN unzip protoc-3.7.1-linux-x86_64.zip -d protoc3
#RUN mv protoc3/bin/* /usr/local/bin/
#RUN mv protoc3/include/* /usr/local/include/google
#RUN chown -R $USER /usr/local/bin/protoc
#RUN chown -R $USER /usr/local/include/google

#ENV PATH=$PATH:/usr/local/bin

#RUN rm -rf protoc3
#RUN rm -rf protoc-3.7.1-linux-x86_64.zip

#RUN go mod vendor

#WORKDIR /go/src/github.com/keti-openfx/openfx/pb
#RUN make 

#WORKDIR /go/src/github.com/keti-openfx/openfx

RUN CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w" \
	-a -installsuffix cgo -o fxgateway .


FROM alpine:3.7

RUN addgroup -S app \
	&& adduser -S -g app app \
	&& apk --no-cache add \
	&& mkdir /etc/docker \
#	&& echo "{ "dns": ["172.17.0.1"] }" >> /etc/docker/daemon.json \
	ca-certificates
WORKDIR /home/app

EXPOSE 10000

COPY --from=builder /go/src/github.com/keti-openfx/openfx .
RUN chown -R app:app ./

USER app

CMD ["./fxgateway"]
