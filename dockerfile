FROM golang:1.20.1-alpine3.17 AS builder
WORKDIR /tmp
COPY . .
RUN apk add make
RUN make

FROM alpine:latest
WORKDIR /app
COPY --from=builder /tmp/releasebot .
COPY --from=builder /tmp/config.json .
CMD ["./releasebot"]
