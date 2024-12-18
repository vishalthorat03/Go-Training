# Start with Golang base image
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Install necessary build tools and dependencies
RUN apk add --no-cache gcc musl-dev make bash

# Copy the Go module files and install dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the application port
EXPOSE 9090

# Run the application
CMD ["./main"]
