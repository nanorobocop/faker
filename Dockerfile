FROM golang:1.17-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o faker .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /src/faker .
CMD [ "/app/faker" ]
