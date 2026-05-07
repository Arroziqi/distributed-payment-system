FROM golang:alpine AS builder
WORKDIR /app
COPY internal/shared internal/shared
COPY services/auth-service/go.mod services/auth-service/go.sum services/auth-service/
WORKDIR /app/services/auth-service
RUN go mod download
COPY services/auth-service .
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/auth-service ./cmd/auth-service

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/auth-service /usr/local/bin/auth-service
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/auth-service"]
