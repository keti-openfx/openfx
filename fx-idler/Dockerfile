FROM golang:1.13-alpine AS builder

WORKDIR /go/src/github.com/keti-openfx/fx-idler

ENV GO111MODULE=on
ENV CGO_ENABLED=0

COPY . . 

RUN gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*")

RUN GOOS=linux go build --ldflags "-s -w" \
	-a -installsuffix cgo -o /usr/bin/fx-idler .

FROM alpine:3.12

RUN addgroup -S app && adduser -S -g app app
RUN mkdir -p /home/app

WORKDIR /home/app

COPY --from=builder /usr/bin/fx-idler /home/app/

RUN chown -R app /home/app
USER app

EXPOSE 8080
VOLUME /tmp

ENTRYPOINT ["/home/app/fx-idler"]
