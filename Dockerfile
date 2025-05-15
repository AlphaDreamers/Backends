RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

# Add CA certificates and timezone data
RUN apk add --no-cache supervisor ca-certificates tzdata

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/main .

# Make sure it's executable
RUN chmod +x /app/main

COPY supervisord.conf /etc/supervisord.conf

RUN mkdir -p /var/log

# Expose all the service ports
EXPOSE 8001 8002 8003 8004 8005 8006

# Environment variables for service ports
ENV AUTH_PORT=8001 \
    CHAT_PORT=8002 \
    WALLET_PORT=8003 \
    ORDER_PORT=8004 \
    GIG_PORT=8005 \
    USER_PORT=8006

# Set the entrypoint to supervisord
ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]