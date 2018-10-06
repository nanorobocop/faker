FROM golang:1.11
WORKDIR /usr/src/app
COPY http-test-server.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/http-test-server http-test-server.go

FROM alpine:latest
WORKDIR /usr/src/app
COPY --from=0 /usr/src/app/bin/http-test-server .
CMD [ "http-test-server" ]