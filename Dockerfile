# ---------- Build stage ----------
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy Go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go application for Linux, statically
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# ---------- Final stage ----------
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the statically compiled binary from builder
COPY --from=builder /app/app .

# Expose the app port
EXPOSE 8001

# Command to run the executable
CMD ["./app"]
