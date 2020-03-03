FROM golang:1.12 as builder

RUN mkdir -p /go/src/github.com/keti-openfx/openfx

WORKDIR /go/src/github.com/keti-openfx/openfx

COPY . .

ENV GO111MODULE=on
RUN gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*")

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
