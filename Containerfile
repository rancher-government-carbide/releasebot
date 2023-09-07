FROM golang:1.20.1-alpine3.17 AS builder
COPY . /build
WORKDIR /build
RUN apk add make git 
RUN apk add --no-cache ca-certificates
RUN adduser -D nonroot
RUN make dependencies
RUN make

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build/releasebot /
COPY payloads.json /
COPY repos.json /
USER nonroot
ENTRYPOINT [ "/releasebot" ]
