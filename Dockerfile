FROM golang:1.15.6-alpine3.12 AS build

WORKDIR /go/src/app

# ARGS does not work outside IMAGE
ARG MODE="dev" 

RUN apk add --no-cache build-base \ 
&& apk add --no-cache wget \ 
&& apk add  --no-cache curl \
&& curl -sfL $(curl -s https://api.github.com/repos/powerman/dockerize/releases/latest | grep -i /dockerize-$(uname -s)-$(uname -m)\" | cut -d\" -f4) | install /dev/stdin /usr/local/bin/dockerize

COPY . .

# If mod arg is equal to dev then rename .env.example to .env else .ev.testing to .env
RUN if [[ ${MODE} == "dev" ]]; then mv .env.example .env ; else mv .env.testing .env ; fi 

RUN go build -o app

# Use alpine image
FROM alpine 

WORKDIR /app

# Here copy our builded app from /go/src/app to /app/
COPY --from=build /go/src/app/app /app/
# Copy ENV
COPY --from=build /go/src/app/.env /app/

EXPOSE 3000

ENTRYPOINT ["./app"]

