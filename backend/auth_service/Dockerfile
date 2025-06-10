# Stage 1 - build
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download && \
    go build -o auth_service cmd/auth_service/main.go
# Stage 2 - run
FROM alpine
WORKDIR /app
COPY --from=builder /app/auth_service .
CMD ["./auth_service"]