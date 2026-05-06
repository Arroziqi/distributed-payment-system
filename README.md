# Distributed Payment System (Golang)

Production-style distributed payment system using microservices and clean architecture.

## Services

- `auth-service`: JWT + refresh token lifecycle
- `wallet-service`: balance, topup, withdrawal, ledger
- `transaction-service`: transfer orchestration, history, idempotency
- `notification-service`: asynchronous notifications and retries

## Infra

- PostgreSQL (database per service)
- Redis (idempotency cache + distributed locking)
- RabbitMQ (event bus)
- Prometheus (metrics scraping)
- Grafana (dashboards)

## Run

```bash
docker compose up --build
```

Endpoints:

- Auth: `http://localhost:8081`
- Wallet: `http://localhost:8082`
- Transaction: `http://localhost:8083`
- Notification: `http://localhost:8084`
- RabbitMQ UI: `http://localhost:15672` (`guest/guest`)
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (`admin/admin`)

## API Contract (v1)

### Auth

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `POST /auth/logout`

### Wallet

- `POST /wallet/topups` (header `Idempotency-Key`)
- `POST /wallet/withdrawals` (header `Idempotency-Key`)
- `GET /wallets/:userId/balance`

### Transaction

- `POST /transactions/transfers` (header `Idempotency-Key`)
- `GET /transactions`
- `GET /transactions/:id`

### Notification

- `GET /notifications/:id`

## Event Flow

- `auth.user.registered` -> wallet creates default wallet
- `wallet.topup.completed` -> notification sends topup confirmation
- `wallet.withdrawal.completed` -> notification sends withdrawal confirmation
- `transaction.transfer.requested` -> wallet performs debit/credit
- `transaction.transfer.completed` -> notification sends transfer success
- `transaction.failed` -> notification sends failure notice

## Concurrency + Idempotency Strategy

- Request-level idempotency key persisted in `idempotency_keys`
- Redis distributed lock per wallet on debit path (`wallet:{id}`)
- SQL transaction with `SELECT ... FOR UPDATE` for balance mutation
- Outbox table written in same SQL transaction to guarantee event durability

## Folder Shape (per service)

```text
cmd/<service>/main.go
internal/domain/
internal/usecase/
internal/repository/
internal/delivery/http/
internal/infrastructure/
migrations/
Dockerfile
go.mod
```
