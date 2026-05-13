# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o smartgorest .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (PostgreSQL client for healthchecks)
RUN apk add --no-cache postgresql-client

# Copy the binary from builder
COPY --from=builder /app/smartgorest .

# Copy public directory for templates and static files
COPY --from=builder /app/public ./public

# Copy .env.example for reference (optional)
COPY --from=builder /app/.env .env.example

# Expose the application port
EXPOSE 1323

# Set the entrypoint
CMD ["./smartgorest"]

