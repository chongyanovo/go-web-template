FROM ubuntu:20.04
LABEL authors="chongyan"

COPY go-web-app /app
WORKDIR /app

ENTRYPOINT ["/app/go-web-app"]