# Stage 1: Build the Go binary
FROM golang:1.17-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o messages-service

# Stage 2: Create the runtime image
FROM alpine:latest

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/messages-service .

# Expose port 8080
EXPOSE 8080

# Command to run when starting the container
CMD ["./messages-service"]