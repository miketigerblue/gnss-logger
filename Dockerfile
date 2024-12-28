# Stage 1: Build
FROM golang:1.20-buster AS builder
WORKDIR /app

# Copy go.mod and go.sum first, then download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o gnss-logger main.go

# Stage 2: Runtime
FROM debian:buster-slim
WORKDIR /app

# Copy from builder
COPY --from=builder /app/gnss-logger /usr/local/bin/gnss-logger

# If you need ca-certificates for TLS etc.:
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Run the logger on container start
ENTRYPOINT ["/usr/local/bin/gnss-logger"]
