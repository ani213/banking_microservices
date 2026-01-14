# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service cmd/auth/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/auth-service .

# Set default environment variable
# ENV EMAIL_SERVER=http://localhost:8083

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./auth-service"]