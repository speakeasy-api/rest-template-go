# syntax=docker/dockerfile:1

# Build the application

FROM golang:1.18-alpine AS build

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN ls -lrt

RUN go build -o ./server ./cmd/server/main.go 

# Build the server image

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /app/server ./
COPY --from=0 /app/config/config-docker.yaml ./config/config.yaml
COPY --from=0 /app/migrations/ ./migrations/

CMD ["./server"]