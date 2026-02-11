FROM ubuntu:latest
LABEL authors="ricardgo"

# Use the official Golang image as the base
FROM golang:1.25-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Download dependencies first (caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the Go binary
RUN ls
RUN go build -o /app ./cmd/main.go

# Run
# Final lightweight stage
FROM alpine:3.23 AS final

# Copy the compiled binary from the builder stage
COPY --from=builder /app /app/crypto
COPY resources/config.yaml resources/config.yaml
EXPOSE 3000
# ENTRYPOINT ["/bin/sh -c"]
CMD ["/app/crypto/main"]

# ENTRYPOINT ["top", "-b"]