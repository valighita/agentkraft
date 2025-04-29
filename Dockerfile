# Use the official Golang image as the base image
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o agentkraft cmd/main.go

# Use a minimal base image for the final container
FROM debian:bookworm-slim

# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/agentkraft .
COPY --from=builder /app/server/frontend/ ./server/frontend

# Command to run the application
CMD ["./agentkraft"]
