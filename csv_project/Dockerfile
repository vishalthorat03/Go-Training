# Build Stage
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application code and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o csv_uploader

# Production Stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Install required packages
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Kolkata

# Create necessary directories and ensure permissions
RUN mkdir -p /app && chmod -R 777 /app
RUN touch /app/app.log && chmod 666 /app/app.log

# Add a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Copy the built binary and frontend files
COPY --from=builder /app/csv_uploader .
COPY frontend ./frontend

# Expose the application port
EXPOSE 4041

# Run the application
CMD ["./csv_uploader"]
