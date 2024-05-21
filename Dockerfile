# Use official Golang image as the base image
FROM golang:1.22-alpine AS builder

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download necessary Go modules
COPY go.mod .
COPY go.sum .
# COPY config/config.yaml ./config/config.yaml
RUN go mod download

# Copy the code into the container
COPY . .

# Build the Go app
RUN go build -o app .

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/app .
COPY config/config.yaml ./config/config.yaml

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./app", "serve", "--config", "config.yaml", "--debug", "true"]
