FROM golang:1.10.1 as builder

RUN mkdir -p /go/src/github.com/keti-openfx/openfx-gateway

WORKDIR /go/src/github.com/keti-openfx/openfx-gateway

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w" \
	-a -installsuffix cgo -o fxgateway .


FROM alpine:3.7

RUN addgroup -S app \
	&& adduser -S -g app app \
	&& apk --no-cache add \
	ca-certificates
WORKDIR /home/app

EXPOSE 10000

COPY --from=builder /go/src/github.com/keti-openfx/openfx-gateway .
RUN chown -R app:app ./

USER app

CMD ["./fxgateway"]
