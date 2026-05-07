# Requirements

## Functional Requirements
- **Authentication**: Users can register and log in. System issues JWT tokens.
- **Wallet Operations**: Users can create a wallet and query their balance.
- **Transfers**: Users can top up their wallet, transfer funds to other wallets, and withdraw funds.
- **Notifications**: Users receive asynchronous notifications regarding their transaction status.
- **Frontend Dashboard**: A comprehensive UI allowing users to execute and view all above operations.

## Non-Functional Requirements
- **Concurrency Safety**: The system must handle concurrent transactions safely (e.g., preventing double-spending) using database row-level locks or Go concurrency mechanisms.
- **Scalability**: Services must be stateless and horizontally scalable.
- **Observability**: All services must expose metrics to Prometheus and push structured logs to Loki.
- **Security**: Secure JWT handling, password hashing, and API authorization.
- **Dockerization**: The entire stack must be easily orchestratable via Docker Compose.
- **Fault Tolerance**: Services should gracefully handle failures of dependencies (e.g., retries for RabbitMQ connection).
