FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/transaction-service ./cmd/transaction-service

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/transaction-service /usr/local/bin/transaction-service
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/transaction-service"]
