# # Build Stage
# FROM golang:1.23 AS builder

# # Set the working directory inside the container
# WORKDIR /app

# # Copy Go modules and download dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the entire application code and build the binary
# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o csv_uploader

# # Set the working directory inside the container
# WORKDIR /app

# # Install required packages
# RUN apk --no-cache add ca-certificates tzdata
# ENV TZ=Asia/Kolkata

# # # Create necessary directories and ensure permissions
# # RUN mkdir -p /app/logs && chmod -R 777 /app/logs
# # RUN touch /app/logs/app.log && chmod 666 /app/logs/app.log

# # Create the log directory and set appropriate permissions
# RUN mkdir -p /app/app/logs/ && chmod -R 777 /app/app/logs

# # Create the log file and set appropriate permissions
# RUN touch /app/app/logs/app.log && chmod 666 /app/app/logs/app.log

# # Add a non-root user for security
# RUN addgroup -S appgroup && adduser -S appuser -G appgroup
# RUN chown -R appuser:appgroup /app

# # Switch to the non-root user
# USER appuser

# # Copy the built binary from the builder stage
# # COPY --from=builder /app/csv_uploader /app/csv_uploader
# COPY --from=builder /app/csv_uploader /app/csv_uploader

# # Copy frontend files (if any)
# COPY frontend /app/frontend

# # Expose the application port
# EXPOSE 4041

# # Run the application
# # CMD ["./csv_uploader"]
# CMD ["/app/csv_uploader"]


# Build Stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application code and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o csv_uploader

# Production Stage
FROM alpine:latest

WORKDIR /app

# Install required packages
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Kolkata

# Ensure the log directory and file exist with correct permissions
RUN mkdir -p /app/app/logs && chmod -R 777 /app/app/logs && touch /app/app/logs/app.log && chmod 666 /app/app/logs/app.log

# Add a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /app

# Switch to root user for debugging log writing issues
USER root

# Copy the built binary from the builder stage
COPY --from=builder /app/csv_uploader /app/csv_uploader

# Expose the application port
EXPOSE 4041

# Run the application
CMD ["/app/csv_uploader"]
