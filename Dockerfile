# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder

WORKDIR /app

# Copy go files
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy everything else
COPY . ./

# Build the app
RUN go build -o urlshortener ./cmd/server

# Final minimal image
FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/urlshortener .

EXPOSE 8080

CMD ["./urlshortener"]
