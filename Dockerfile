# Stage 1: Build the Go binary
FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o tiny ./cmd/main.go

EXPOSE 42000

# Run the app
CMD ["./tiny"]
