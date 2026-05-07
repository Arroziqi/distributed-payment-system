FROM golang:alpine AS builder
WORKDIR /app
COPY internal/shared internal/shared
COPY services/transaction-service/go.mod services/transaction-service/go.sum services/transaction-service/
WORKDIR /app/services/transaction-service
RUN go mod download
COPY services/transaction-service .
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/transaction-service ./cmd/transaction-service

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/transaction-service /usr/local/bin/transaction-service
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/transaction-service"]
