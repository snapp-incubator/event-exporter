#build stage
FROM golang:1.19.1-buster AS builder

WORKDIR /go/src/app

COPY go.sum go.mod /go/src/app/

RUN go mod download

COPY . /go/src/app
RUN go build -ldflags="-w -s" -o event_exporter

#final stage
FROM debian:buster-slim

ENV TZ=UTC \
    PATH="/app:${PATH}"

RUN mkdir -p /var/log && \
    chgrp -R 0 /var/log && \
    chmod -R g=u /var/log
WORKDIR /app

COPY --from=builder /go/src/app/event_exporter /app/event_exporter
ENTRYPOINT /app/event_exporter
LABEL Name=event_exporter Version=1.0.0
EXPOSE 8090
