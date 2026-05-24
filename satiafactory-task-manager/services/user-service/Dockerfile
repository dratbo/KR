FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/user-service .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8081
CMD ["./user-service"]