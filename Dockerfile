# =========================
# Build stage
# =========================
FROM docker.io/library/golang:1.25.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/server

# =========================
# Runtime stage
# =========================
FROM docker.io/library/alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    imagemagick \
    poppler-utils \
    ca-certificates \
    tzdata \
    ghostscript \
    ghostscript-fonts

# Create non-root user
RUN addgroup -S app && adduser -S app -G app

# Application directory (read-only)
WORKDIR /app
COPY --from=builder /app/app /app/app

# Writable data directory
RUN mkdir -p /storage/uploads /storage/cache /storage/tmp \
    && chown -R app:app /storage \
    && chmod -R 755 /storage

# Environment variables
ENV TZ=Asia/Jakarta

# Drop privileges
USER app

EXPOSE 8080

CMD ["/app/app"]
