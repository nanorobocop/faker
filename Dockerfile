FROM golang:1.17-alpine AS builder
WORKDIR /src
COPY faker.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o faker faker.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /src/faker .
CMD [ "/app/faker" ]
