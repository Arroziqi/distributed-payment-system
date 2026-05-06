# Distributed Payment System (Golang)

A production-ready distributed payment system built using microservices and Clean Architecture principles in Go.

## Architecture Overview

The system consists of four independent microservices that communicate synchronously via HTTP and asynchronously via RabbitMQ. Each service manages its own PostgreSQL database to ensure loose coupling.

1. **Auth Service**: Manages user registration, authentication, and JWT lifecycle (Access/Refresh tokens).
2. **Wallet Service**: Manages user balances. Uses row-level locking (`SELECT ... FOR UPDATE`) to prevent race anomalies during concurrent top-ups, withdrawals, and transfers.
3. **Transaction Service**: Orchestrates money transfers between wallets. Employs idempotency keys to prevent duplicate transactions and publishes events to the message broker upon success.
4. **Notification Service**: Consumes events from RabbitMQ and handles asynchronous notifications to users.

### Tech Stack
- **Language**: Go (Golang)
- **Databases**: PostgreSQL (per service), Redis (caching and idempotency)
- **Message Broker**: RabbitMQ
- **Observability**: Prometheus (Metrics), Grafana (Dashboards)
- **Documentation**: Swagger/OpenAPI

## Setup Instructions & Docker Compose Usage

The entire environment is containerized using Docker and Docker Compose. You do not need to install local databases or message brokers.

To start the system:

```bash
# Clone the repository
# Navigate to the project root

# Start all services and dependencies in detached mode
docker compose up --build -d

# Check the logs of all services
docker compose logs -f
```

To shut down the system and remove volumes:
```bash
docker compose down -v
```

## Service Ports & URLs

Once the system is running via `docker compose`, the services are bound to the following local ports:

| Service | Port | Base URL | Swagger UI URL |
|---------|------|----------|----------------|
| **Auth** | `8081` | `http://localhost:8081` | [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html) |
| **Wallet** | `8082` | `http://localhost:8082` | [http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html) |
| **Transaction** | `8083` | `http://localhost:8083` | [http://localhost:8083/swagger/index.html](http://localhost:8083/swagger/index.html) |
| **Notification** | `8084` | `http://localhost:8084` | [http://localhost:8084/swagger/index.html](http://localhost:8084/swagger/index.html) |

**Other Infrastructure:**
- RabbitMQ Management UI: `http://localhost:15672` (Creds: `guest` / `guest`)
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (Creds: `admin` / `admin`)

## Public vs Protected Endpoints

**Public Endpoints (No Auth Required):**
- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `GET /healthz` (on all services)
- `GET /metrics` (on all services)

**Protected Endpoints (Requires Bearer Token):**
All internal application logic requires authentication via a JWT Bearer Token, including:
- All **Wallet** endpoints (Top up, Withdraw, Transfer, Balance Check)
- All **Transaction** endpoints (Process Payment, List History, Detail)
- All **Notification** endpoints
- `POST /auth/logout`

*Note: You can easily test protected endpoints using the **Authorize** button in each service's Swagger UI.*

## API Testing Flow Summary

For a full step-by-step tutorial on testing the endpoints using Swagger, see [API_USAGE.md](./API_USAGE.md).

**Quick Summary:**
1. **Register** a user via Auth Service (`POST /auth/register`).
2. **Login** to retrieve an Access Token (`POST /auth/login`).
3. **Authorize** in Swagger by pasting the Access Token in the `Authorize` dialog (`Bearer <token>`).
4. **Create a Wallet** in Wallet Service (`POST /wallets`).
5. **Top Up** your balance (`POST /wallet/topups`).
6. **Transfer** funds via Transaction Service (`POST /transactions/payments`). Make sure to provide an `Idempotency-Key` header.
7. **Check History** using Transaction Service (`GET /transactions`).

For visual sequence diagrams of these flows, check out [SYSTEM_FLOW.md](./SYSTEM_FLOW.md).

## Event Flow & Concurrency

- `transaction.transfer.completed` -> Published to RabbitMQ -> Notification Service sends transfer success.
- **Concurrency**: Handled via SQL `SELECT ... FOR UPDATE` ensuring safe balance mutation.
- **Idempotency**: Cached in Redis to guarantee retried requests (like flaky network transfers) don't duplicate ledger entries.
