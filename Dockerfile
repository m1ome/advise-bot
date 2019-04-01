# Building application
FROM golang:1.10-alpine3.8 as builder

COPY . /go/src/github.com/m1ome/advise-bot
WORKDIR /go/src/github.com/m1ome/advise-bot

RUN apk add --no-cache git gcc musl-dev
RUN go get -u github.com/golang/dep/...
RUN dep ensure

RUN \
    export VERSION=$(git rev-parse --verify HEAD) && \
    export LDFLAGS="-w -s -X main.Version=${VERSION}" && \
    export CGO_ENABLED=0 && \
    go build -v -ldflags "${LDFLAGS}" -o /go/bin/bot .

# Packing in main container
FROM alpine:3.8

RUN apk add --no-cache ca-certificates

WORKDIR /

COPY --from=builder /go/bin/bot /bot

CMD ["/bot"]