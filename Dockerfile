FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api
# RUN CGO_ENABLED=0 GOOS=linux go build -o seeder ./cmd/seeder

FROM alpine:latest

WORKDIR /root
RUN apk add --no-cache ca-certificates wget

COPY --from=builder /app/main .
# COPY --from=builder /app/seeder .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

CMD ["./main"]
