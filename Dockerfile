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

# Build the application
RUN go build -o main ./cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Add necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy config files
COPY --from=builder /app/.env* ./
COPY --from=builder /app/config* ./

# Create necessary directories if needed
RUN mkdir -p /app/logs

# Set executable permissions
RUN chmod +x /app/main

# Expose the port
EXPOSE 3000

# Command to run the application
CMD ["/app/main"] 