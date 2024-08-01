# Build stage: Build the Go application
FROM golang:1.22-alpine AS builder

# Install git
RUN apk update && apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Set the working directory to the location of main.go
WORKDIR /app/cmd/weblayout

# Build the Go app
RUN go build -o /app/main .

# Deploy stage: Create a minimal runtime image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .

# Copy the .env file from the build context to the runtime image
COPY configs/config.env /app/configs/config.env
# Set the CONFIG_PATH environment variable
ENV CONFIG_PATH=/app/configs/config.env

# Command to run the executable
CMD ["./main"]
