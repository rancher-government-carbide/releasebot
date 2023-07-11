FROM golang:1.20.1-alpine3.17 AS builder
COPY . /build
WORKDIR /build
RUN apk add make git 
RUN apk add --no-cache ca-certificates
RUN make dependencies
RUN make

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/releasebot /
COPY config.json /
CMD ["/releasebot"]
