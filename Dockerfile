FROM golang:1.11.4-alpine3.8 AS build
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    http_proxy=http://snapp-mirror:TmfBZb68qjGGF6feBdqX@mirror-fra-1.snappcloud.io:30128 \
    https_proxy=http://snapp-mirror:TmfBZb68qjGGF6feBdqX@mirror-fra-1.snappcloud.io:30128
RUN apk add --update --no-cache \
      build-base \
      git \
      ca-certificates \
    && \
    mkdir -p /src
COPY go.sum go.mod /src/
WORKDIR /src
RUN go mod download
COPY . /src
RUN go build -ldflags="-w -s" -o event-exporter

FROM alpine:3.8
ENV TZ=UTC \
    OPENSHIFT_NAMESPACE=default
RUN apk add --update --no-cache \
      tzdata \
      ca-certificates curl \
    && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && unset http_proxy && unset https_proxy && \
    echo "${TZ}" > /etc/timezone
WORKDIR /app
EXPOSE 8090
COPY --from=build /src/event-exporter /app/
CMD ["/app/event-exporter"]
USER 1001