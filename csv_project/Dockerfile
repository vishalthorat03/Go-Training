# Stage 1: Build the Go application
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Install the required dependencies for your app
RUN go get github.com/shirou/gopsutil/cpu
RUN go get github.com/shirou/gopsutil/mem
RUN go get github.com/jackc/pgx/v4
RUN go get github.com/stretchr/testify/assert

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o csv_uploader

# Stage 2: Use a lightweight Alpine image for runtime
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

# Install build dependencies for MockGain
RUN apk add --no-cache git make gcc

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/csv_uploader .
COPY frontend ./frontend

# Create a directory for file uploads
RUN mkdir -p /app/uploads

# Ensure the log file is created with appropriate permissions (optional)
RUN touch /app/app.log && chmod 666 /app/app.log

# Debug: List the files in the directory
RUN ls -l /app

EXPOSE 4041

# Run the Go application
CMD ["./csv_uploader"]
