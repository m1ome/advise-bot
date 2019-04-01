# Building application
FROM golang:1.10-alpine3.8 as builder

COPY . /go/src/github.com/m1ome/advise-bot
WORKDIR /go/src/github.com/m1ome/advise-bot

RUN apk add --no-cache git gcc musl-dev

# Packing in main container
FROM alpine:3.8

RUN apk add --no-cache ca-certificates

WORKDIR /

COPY --from=builder /go/bin/bot /bot

CMD ["/bot"]