FROM golang:alpine AS builder
WORKDIR /app
COPY internal/shared internal/shared
COPY services/wallet-service/go.mod services/wallet-service/go.sum services/wallet-service/
WORKDIR /app/services/wallet-service
RUN go mod download
COPY services/wallet-service .
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/wallet-service ./cmd/wallet-service

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/wallet-service /usr/local/bin/wallet-service
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/wallet-service"]
