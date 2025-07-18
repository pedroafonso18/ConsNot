# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o robo_consultas ./cmd/ConsNot

# --- Run stage ---
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/robo_consultas ./robo_consultas
COPY .env ./
# If you need CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Set environment variables if needed
# ENV VAR_NAME=value

ENTRYPOINT ["./robo_consultas"] 