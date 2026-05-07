# Architecture

## High-Level Architecture
The system consists of a Vue 3 frontend communicating with multiple Golang microservices. The backend services utilize PostgreSQL for persistence, Redis for caching, and RabbitMQ for asynchronous event-driven communication. An observability stack monitors the ecosystem.

### Frontend
- **Vue 3 + Quasar**: Atomic Design structure with Storybook. Interacts with backend services via Axios.

### Backend Services
1. **Auth Service**: Manages user registration, login, and JWT token issuance.
2. **Wallet Service**: Handles wallet creation and balance queries.
3. **Transaction Service**: Processes top-ups, transfers, and withdrawals. Enforces concurrency control.
4. **Notification Service**: Consumes RabbitMQ events and stores notifications for the user.

### Infrastructure & Data Stores
- **PostgreSQL**: Primary relational database. Each service has its own database (auth_db, wallet_db, transaction_db, notification_db) to enforce data boundaries.
- **Redis**: In-memory data store for caching and rate-limiting.
- **RabbitMQ**: Message broker for asynchronous inter-service communication.

### Observability Stack
- **Prometheus**: Scrapes metrics from backend services.
- **Grafana**: Visualizes metrics and logs.
- **Loki**: Centralized log aggregation system.
- **Promtail**: Ships logs from Docker containers to Loki.

## Communication Flow
1. **Synchronous (HTTP/REST)**: Frontend -> Backend Services. Also, potential API Gateway -> Backend Services.
2. **Asynchronous (AMQP)**: Transaction Service -> RabbitMQ -> Notification Service.
