FROM golang:1.17.2-alpine AS build

RUN apk add --no-cache build-base wget curl \
&& curl -sfL $(curl -s https://api.github.com/repos/powerman/dockerize/releases/latest | grep -i /dockerize-$(uname -s)-$(uname -m)\" | cut -d\" -f4) | install /dev/stdin /usr/local/bin/dockerize

WORKDIR /go/src/app

COPY go.mod go.sum .

RUN go mod download

COPY . .

ARG MODE="dev"

RUN set -ex; \
  if [[ ${MODE} == "dev" ]]; then mv .env.example .env; \
  elif [[ ${MODE} == "docker" ]]; then mv .env.docker .env ; \
  else mv .env.testing .env; fi; \
  mkdir /app; \
  cp .env /app/.env

RUN go build -o /app/app

FROM alpine
WORKDIR /app

COPY --from=build /usr/local/bin/dockerize /usr/local/bin/dockerize
COPY --from=build /app/ ./

EXPOSE 3000
ENTRYPOINT ["/app/app"]
