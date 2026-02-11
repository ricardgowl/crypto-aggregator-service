FROM ubuntu:latest
LABEL authors="ricardgo"

# Build
FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/api ./cmd/api

# Run
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /bin/api /api
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/api"]

# ENTRYPOINT ["top", "-b"]