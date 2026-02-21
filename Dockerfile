FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install semua C dependencies dulu sebelum build apapun
RUN apk add --no-cache \
    git \
    gcc \
    musl-dev \
    tesseract-ocr-dev \
    leptonica-dev \
    pkgconfig

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go --output docs

# Build API (tidak butuh CGO)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Build BOT (butuh CGO untuk gosseract)
RUN CGO_ENABLED=1 GOOS=linux go build -o bot ./cmd/bot

# ------------------------------------------------
FROM alpine:latest

WORKDIR /root

# Install runtime dependencies (tesseract + bahasa)
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    tesseract-ocr \
    tesseract-ocr-data-ind \
    tesseract-ocr-data-eng \
    leptonica \
    libstdc++ \
    libgcc

ENV TZ=Asia/Jakarta
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

COPY --from=builder /app/main .
COPY --from=builder /app/bot .
COPY --from=builder /app/docs ./docs

EXPOSE 8080