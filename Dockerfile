# syntax = docker/dockerfile:latest
FROM golang:1.19.1-alpine3.16 AS build

ARG GOPATH="/go"
ARG GOMODCACHE=${GOPATH}/pkg/mod
ARG GOCACHE="/root/.cache/go-build"
ARG MODE="dev"

ENV GOMODCACHE=${GOMODCACHE}
ENV GOCACHE=${GOCACHE}

RUN set -ex; \
    apk update && \
    apk add --no-cache build-base wget curl;

COPY go.* ./

RUN --mount=type=cache,mode=0777,target=${GOCACHE} \
    --mount=type=cache,mode=0777,target=${GOMODCACHE} \
    go mod download

COPY . ./

RUN --mount=type=cache,mode=0777,target=${GOCACHE} \
    --mount=type=cache,mode=0777,target=${GOMODCACHE} \
    set -ex; \
    if [[ ${MODE} == "dev" ]]; then mv .env.example .env; \
    elif [[ ${MODE} == "docker" ]]; then mv .env.docker .env ; \
    else mv .env.testing .env; fi; \
    mkdir /app; \
    cp -r assets /app/assets; \
    cp -r config /app/config; \
    cp -r database /app/database; \
    cp .env /app/.env; \
    go build -o /tmp/app; \
    cp -r /tmp/app /app/;

FROM alpine 
WORKDIR /app

COPY --from=build /tmp/app /app/ ./

ENTRYPOINT ["./app"]
