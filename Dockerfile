# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app

# Enable Go modules & speed up builds with minimal deps layers
COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o toy-api ./cmd/server

FROM alpine:3.19
WORKDIR /app

# Add non-root user for better security
RUN addgroup -S toy && adduser -S -G toy toy

COPY --from=builder /app/toy-api /app/toy-api

EXPOSE 8080
USER toy

ENV PORT=8080

ENTRYPOINT ["/app/toy-api"]
