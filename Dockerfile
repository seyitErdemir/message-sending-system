# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application for production
RUN go build -o main ./cmd/api/main.go

# Development stage
FROM golang:1.21-alpine AS development

WORKDIR /app

# Install necessary development tools
RUN apk add --no-cache git curl

# Install Air for hot reload
RUN go install github.com/cosmtrek/air@v1.49.0

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Set development command
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:latest AS production

WORKDIR /app

# Add necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy config files
COPY --from=builder /app/.env* ./
COPY --from=builder /app/config* ./

# Create necessary directories
RUN mkdir -p /app/logs

# Set executable permissions
RUN chmod +x /app/main

# Expose the port
EXPOSE 3000

# Command to run the application
CMD ["/app/main"] 