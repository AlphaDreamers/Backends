FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/main.go


FROM alpine:latest

RUN apk add --no-cache nginx supervisor

WORKDIR /app

COPY --from=builder /app/main .

COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY supervisord.conf /etc/supervisord.conf

RUN mkdir -p /run/nginx && \
    touch /run/nginx/nginx.pid && \
    mkdir -p /var/log/nginx

EXPOSE 80

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
