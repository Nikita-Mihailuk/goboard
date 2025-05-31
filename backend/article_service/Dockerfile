# Stage 1 - build
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download && \
    go build -o article_service cmd/article_service/main.go
# Stage 2 - run
FROM alpine
WORKDIR /app
COPY --from=builder /app/article_service .
CMD ["./article_service"]