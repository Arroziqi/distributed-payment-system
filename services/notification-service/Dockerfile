FROM golang:alpine AS builder
WORKDIR /app
COPY internal/shared internal/shared
COPY services/notification-service/go.mod services/notification-service/go.sum services/notification-service/
WORKDIR /app/services/notification-service
RUN go mod download
COPY services/notification-service .
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/notification-service ./cmd/notification-service

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/notification-service /usr/local/bin/notification-service
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/notification-service"]
