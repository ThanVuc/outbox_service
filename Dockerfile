# build stage
FROM golang:alpine AS builder
RUN apk add --no-cache

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o outbox_service ./main.go

# stage 2
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/outbox_service .
RUN chmod +x /app/outbox_service
ENTRYPOINT ["./outbox_service"]
