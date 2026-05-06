FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/notification-service ./cmd/notification-service

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/notification-service /usr/local/bin/notification-service
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/notification-service"]
