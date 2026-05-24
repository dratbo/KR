FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o task-service ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/task-service .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8082
CMD ["./task-service"]