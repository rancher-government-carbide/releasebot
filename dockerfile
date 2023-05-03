FROM golang:1.20.1-alpine3.17 AS builder

COPY . /build

WORKDIR /build

ENV GOOS=linux
ENV CGO_ENABLED=0

RUN go get -d -v ./...

RUN go build -v -o releasebot ./cmd

FROM alpine

WORKDIR /app
COPY --from=builder /build/releasebot /usr/local/bin/
COPY config.json .

CMD ["releasebot"]
