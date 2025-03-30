# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy source code
COPY ./weather-app .

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/weather-app .

# Add a non-root user
RUN adduser -D appuser

# Sets environment variable to indicate that the app is running in a Docker container
ENV DOCKER_CONTAINER=true

USER appuser

# Command to run the executable
ENTRYPOINT ["./weather-app"]