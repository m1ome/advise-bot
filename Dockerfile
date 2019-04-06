# Building application
FROM golang:1.12-alpine3.9 as builder

COPY . /src/bot
WORKDIR /src/bot

RUN apk add git

RUN \
    export VERSION=$(git rev-parse --verify HEAD) && \
    export LDFLAGS="-w -s -X main.Version=${VERSION}" && \
    export CGO_ENABLED=0 && \
    go build -v -ldflags "${LDFLAGS}" -o /go/bin/bot .

# Packing in main container
FROM alpine:3.9

RUN apk add ca-certificates

WORKDIR /

COPY --from=builder /go/bin/bot /bot

CMD ["/bot"]