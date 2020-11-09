#build stage
FROM golang:1.14-alpine AS builder
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
RUN apk add --update --no-cache \
    ca-certificates git

COPY go.sum go.mod /go/src/app/
RUN go env -w GOPROXY="https://repo.snapp.tech/repository/goproxy/" && \
    go env -w GONOSUMDB="gitlab.snapp.ir" && \
    go env -w GO111MODULE="on"

RUN go mod download

COPY . /go/src/app
RUN go build -ldflags="-w -s" -o event_exporter

#final stage
FROM alpine:3.12

ENV TZ=UTC \
    PATH="/app:${PATH}"

RUN apk add --update --no-cache \
      tzdata \
      ca-certificates \
    && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone && \
    mkdir -p /var/log && \
    chgrp -R 0 /var/log && \
    chmod -R g=u /var/log
WORKDIR /app

COPY --from=builder /go/src/app/event_exporter /app/event_exporter
ENTRYPOINT /app/event_exporter
LABEL Name=event_exporter Version=1.0.0
EXPOSE 9876
