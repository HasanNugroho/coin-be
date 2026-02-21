FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go --output docs

# build API
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# build BOT
RUN apk add --no-cache tesseract-ocr-dev leptonica-dev
RUN CGO_ENABLED=1 GOOS=linux go build -o bot ./cmd/bot

FROM alpine:latest

WORKDIR /root
RUN apk add --no-cache ca-certificates wget tzdata

ENV TZ=Asia/Jakarta
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

COPY --from=builder /app/main .
COPY --from=builder /app/bot .
COPY --from=builder /app/docs .

EXPOSE 8080