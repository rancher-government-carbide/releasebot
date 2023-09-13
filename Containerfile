FROM golang:1.20.1-alpine3.17 AS builder
COPY . /build
WORKDIR /build
RUN apk add make git 
RUN apk add --no-cache ca-certificates
RUN make dependencies
RUN make
WORKDIR /permissions
RUN echo "releasebot:x:1001:1001::/:" > passwd && echo "releasebot:x:2000:releasebot" > group

FROM scratch
COPY --from=builder /permissions/passwd /etc/passwd
COPY --from=builder /permissions/group /etc/group
USER releasebot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/releasebot /
COPY payloads.json /
COPY repos.json /
ENTRYPOINT [ "/releasebot" ]
