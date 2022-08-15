# syntax=docker/dockerfile:1

# Build a golang image based on https://docs.docker.com/language/golang/build-images

FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./cmd/server/main.go ./cmd/server/main.go
COPY ./internal/ ./internal/

RUN go build -o ./server ./cmd/server/main.go

# Build the server image

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /app/server ./
COPY ./config/ ./config/
COPY ./migrations/ ./migrations/

EXPOSE 8080

CMD ["./server"]