# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o web-analyzer ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/web-analyzer .
# Copy any config files if needed
COPY --from=builder /app/config/ ./config/

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./web-analyzer"]