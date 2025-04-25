FROM golang:1.23.5 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/internal/config/db/migrations /app/internal/config/db/migrations
COPY .env .

EXPOSE 8081
CMD ["./main"]
