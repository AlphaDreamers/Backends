FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8001 8002 8003

# Run the compiled binary
CMD ["./main"]
