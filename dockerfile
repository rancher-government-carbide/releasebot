FROM golang:1.20.1-alpine3.17 AS builder

COPY . /build

WORKDIR /build

ENV GOOS=linux
ENV CGO_ENABLED=0

RUN go get -d -v ./...

RUN apk add make git && \ 
    make

FROM registry.suse.com/bci/bci-micro:15.4

RUN mkdir -p /home/stigatron && \
    chown -R 1000:1000 /home/stigatron && \
    echo "stigatron:x:1000:1000:stigatron:/tmp:/bin/bash" >> /etc/passwd && \
    echo "stigatron:x:1000:" >> /etc/group

USER 1000

WORKDIR /app
COPY --from=builder /build/releasebot /usr/local/bin/
COPY config.json .

CMD ["releasebot"]
