# # Stage 1: Build the Go application
# FROM golang:1.23 AS builder

# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .
# RUN go build -o csv_uploader

# # Stage 2: Use a lightweight image
# FROM centos:8

# WORKDIR /app
# RUN yum install -y ca-certificates && update-ca-trust

# COPY --from=builder /app/csv_uploader .
# COPY frontend ./frontend
# RUN mkdir uploads

# EXPOSE 4041
# CMD ["./csv_uploader"]


# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o csv_uploader

# Stage 2: Use a lightweight Alpine image for runtime
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/csv_uploader .
COPY frontend ./frontend

# Create a directory for file uploads
RUN mkdir uploads

EXPOSE 4041

# Run the Go application
CMD ["./csv_uploader"]
