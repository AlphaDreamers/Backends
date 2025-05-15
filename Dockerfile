# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

# Add CA certificates and timezone data
RUN apk add --no-cache nginx supervisor ca-certificates tzdata

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/main .

# Make sure it's executable
RUN chmod +x /app/main

COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY supervisord.conf /etc/supervisord.conf

RUN mkdir -p /run/nginx && \
    touch /run/nginx/nginx.pid && \
    mkdir -p /var/log/nginx && \
    mkdir -p /var/log

EXPOSE 80

# Set the entrypoint to supervisord
ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]