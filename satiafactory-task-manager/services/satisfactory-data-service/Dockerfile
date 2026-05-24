FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o data-service ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/data-service .
COPY --from=builder /app/migrations ./migrations
COPY data/Docs.json ./data/Docs.json
EXPOSE 8083
CMD ["./data-service"]