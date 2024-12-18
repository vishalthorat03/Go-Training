# Use an official Golang image as a base
FROM golang:1.23-alpine

# Set the current working directory in the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o app .

# Expose port 8080
EXPOSE 8080

# Run the Go application
CMD ["./app"]
