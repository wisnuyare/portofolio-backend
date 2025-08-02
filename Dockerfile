# Multi-stage build for production deployment
# Build stage - Use official Go image to build the application
FROM golang:1.23-alpine AS builder

# Set necessary environment variables for build optimization
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install git for dependency resolution
RUN apk add --no-cache git ca-certificates tzdata

# Create build directory
WORKDIR /build

# Copy and download dependencies (for better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN go build -ldflags="-s -w" -o portfolio-backend cmd/api/main.go

# Production stage - Use minimal alpine image
FROM alpine:3.19

# Install ca-certificates for HTTPS requests and tzdata for timezone
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy built binary from builder stage
COPY --from=builder /build/portfolio-backend .

# Copy migration files (needed for potential schema updates)
COPY --from=builder /build/migrations ./migrations

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port (Cloud Run will set PORT environment variable)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/v1/health || exit 1

# Run the application
CMD ["./portfolio-backend"]