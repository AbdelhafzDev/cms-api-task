# ---------------------
# Development stage
# ---------------------
FROM golang:1.25.7-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git openssl && \
    go install github.com/air-verse/air@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate dev JWT keys
RUN mkdir -p /app/keys && \
    openssl genpkey -algorithm RSA -out /app/keys/private.pem -pkeyopt rsa_keygen_bits:2048 && \
    openssl rsa -in /app/keys/private.pem -pubout -out /app/keys/public.pem

EXPOSE 8080 9090

# Use air for hot reloading in development
CMD ["air", "-c", ".air.toml"]



# ---------------------
# Build stage
# ---------------------
FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cms-api ./cmd/main.go


# ---------------------
# Production stage
# ---------------------
FROM alpine:3.21 AS production

RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/cms-api .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

EXPOSE 8080 9090

CMD ["./cms-api"]
