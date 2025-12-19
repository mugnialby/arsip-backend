FROM docker.io/library/golang:1.25.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/server

FROM docker.io/library/alpine:3.19

WORKDIR /app

RUN apk add --no-cache \
    imagemagick \
    poppler-utils \
    ca-certificates \
    tzdata

COPY --from=builder /app/app /app/app

RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

CMD ["/app/app"]
