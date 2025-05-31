# Stage 1 - build
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o api_gateway cmd/gateway/main.go

# Stage 2 - run
FROM alpine
WORKDIR /app
COPY --from=builder /app/api_gateway .
CMD ["./api_gateway"]