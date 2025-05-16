FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/main.go

# Use a smaller runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8001
CMD ["./app"]