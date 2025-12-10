# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Download dependencies first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary and migrations
COPY --from=builder /build/server .
COPY --from=builder /build/migrations ./migrations

# Create data directory for SQLite
RUN mkdir -p /app/data

EXPOSE 3000

CMD ["./server"]
