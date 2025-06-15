# Stage 1 - build, использование образа на ubuntu для поддержки работы librdkafka
FROM golang:1.24.2 AS builder

WORKDIR /app

COPY . .

# скачивание librdkafka для confluent-kafka-go
RUN apt-get update && \
    apt-get install -y gcc librdkafka-dev

RUN go mod download && \
    go build -o comment_service cmd/comment_service/main.go

# Stage 2 - run
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/comment_service .
CMD ["./comment_service"]