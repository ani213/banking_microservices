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
RUN CGO_ENABLED=0 GOOS=linux go build -o account-service cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/account-service .

# Set default environment variable
# ENV EMAIL_SERVER=http://localhost:8083 \
#     DATABASE_URL=postgres://postgres:password@localhost:5432/bank?sslmode=disable \
#     SECRET_KEY=super-secret-key

# Expose the port the app runs on
EXPOSE 8081

# Command to run the application
CMD ["./account-service"]