# Distributed Payment System (Golang)

A production-ready distributed payment system built using microservices and Clean Architecture principles in Go.

## Architecture Overview

The system consists of four independent microservices that communicate synchronously via HTTP and asynchronously via RabbitMQ. Each service manages its own PostgreSQL database to ensure loose coupling.

1. **Auth Service**: Manages user registration, authentication, profile management, and JWT lifecycle (Access/Refresh tokens).
2. **Wallet Service**: Manages user balances. Uses **Optimistic Locking** (versioning) to prevent race anomalies during concurrent top-ups, withdrawals, and transfers, ensuring high throughput without blocking.
3. **Transaction Service**: Orchestrates money transfers between wallets. Employs **Pessimistic Locking** for idempotency management and publishes events to the message broker upon success.
4. **Notification Service**: Consumes events from RabbitMQ and handles asynchronous notifications to users.

## Tech Stack Overview

Our system leverages a modern, robust, and scalable stack designed for high performance and maintainability. For a detailed view of our project milestones, see the [Versioning Roadmap](./docs/ROADMAP.md).

### Core Technologies
- **[Go (Golang)](https://go.dev/)**: The backbone of our microservices. Chosen for its exceptional performance, efficient concurrency model (Goroutines), and strong typing.
- **[Vue 3](https://vuejs.org/)**: A progressive JavaScript framework used for building the frontend dashboard. We utilize the **Composition API** and **Atomic Design** for maximum scalability.
- **[Pinia](https://pinia.vuejs.org/)**: The intuitive, type-safe state management store for Vue, handling centralized authentication and user state.
- **[Vite](https://vitejs.dev/)**: A lightning-fast build tool that powers our frontend development experience.
- **[Tailwind CSS](https://tailwindcss.com/)**: A utility-first CSS framework for rapid UI development.
- **[Shadcn-vue](https://www.shadcn-vue.com/)**: A collection of high-quality, accessible UI components.
- **[Orval](https://orval.dev/)**: OpenAPI to TypeScript generator, ensuring type-safe API consumption from our Go backends.

### Data & Messaging
- **[PostgreSQL](https://www.postgresql.org/)**: Our primary relational database. Each microservice manages its own dedicated instance to ensure strict data isolation.
- **[Redis](https://redis.io/)**: Used for high-speed caching and ensuring transaction **idempotency**, preventing duplicate processing.
- **[RabbitMQ](https://www.rabbitmq.com/)**: A robust message broker for asynchronous communication (e.g., triggering notifications).

### Infrastructure & Observability
- **[Docker & Docker Compose](https://www.docker.com/)**: Containerization technology for consistent environments.
- **[Loki & Promtail](https://grafana.com/oss/loki/)**: Centralized logging solution for aggregating logs from all containers.
- **[Prometheus](https://prometheus.io/)**: Monitoring system for collecting real-time metrics.
- **[Grafana](https://grafana.com/)**: Visualization layer for metrics and logs with beautiful, real-time dashboards.

### Tooling & Documentation
- **[Swagger/OpenAPI](https://swagger.io/)**: Automatically generated API documentation for all services.
- **[golang-migrate](https://github.com/golang-migrate/migrate)**: Handles database schema versioning and migrations.
- **[Storybook](https://storybook.js.org/)**: Used for developing and documenting frontend UI components in isolation.

## Setup Instructions & Docker Compose Usage

The entire environment is containerized. To start the system:

```bash
# Clone the repository and navigate to the project root
docker compose up --build -d

# Check the logs
docker compose logs -f
```

To shut down:
```bash
docker compose down -v
```

## Service Ports & URLs

Once running, services are bound to the following local ports:

| Service | Port | Base URL | Swagger UI URL / Notes |
|---------|------|----------|----------------|
| **Frontend UI** | `8000` | `http://localhost:8000` | Vue 3 + Tailwind CSS App |
| **Storybook** | `6006` | `http://localhost:6006` | Component Documentation |
| **Auth** | `8081` | `http://localhost:8081` | [/swagger/index.html](http://localhost:8081/swagger/index.html) |
| **Wallet** | `8082` | `http://localhost:8082` | [/swagger/index.html](http://localhost:8082/swagger/index.html) |
| **Transaction** | `8083` | `http://localhost:8083` | [/swagger/index.html](http://localhost:8083/swagger/index.html) |
| **Notification** | `8084` | `http://localhost:8084` | [/swagger/index.html](http://localhost:8084/swagger/index.html) |

**Infrastructure:**
- RabbitMQ: `http://localhost:15672` (guest/guest)
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (admin/admin)

## Public vs Protected Endpoints

**Public Endpoints (No Auth Required):**
- `POST /auth/register` - User registration
- `POST /auth/login` - Login to get JWT
- `POST /auth/refresh` - Refresh access token
- `GET /healthz` & `GET /metrics` (on all services)

**Protected Endpoints (Requires Bearer Token):**
- **Auth**: `GET /auth/me` (Profile), `PUT /auth/me` (Update), `POST /auth/logout`
- **Wallet**: `POST /wallets` (Create), `POST /wallet/topups`, `POST /wallet/withdrawals`, `POST /wallet/transfers`, `GET /wallets/:userID/balance`
- **Transaction**: `POST /transactions/payments`, `GET /transactions` (History), `GET /transactions/:id` (Detail)

## Full Documentation

Detailed documentation is available in the `/docs` directory:

| File | Description |
|------|-------------|
| [ARCHITECTURE.md](./docs/ARCHITECTURE.md) | Detailed Clean Architecture and microservices design. |
| [API_USAGE.md](./docs/API_USAGE.md) | Step-by-step guide for testing the API flow. |
| [SEQUENCE_DIAGRAMS.md](./docs/SEQUENCE_DIAGRAMS.md) | Mermaid diagrams for system interactions. |
| [SYSTEM_FLOW.md](./docs/SYSTEM_FLOW.md) | Backend logic and event-driven architecture details. |
| [FRONTEND_ARCHITECTURE.md](./docs/FRONTEND_ARCHITECTURE.md) | Vue 3, Pinia, and Vite setup details. |
| [FRONTEND_FLOW.md](./docs/FRONTEND_FLOW.md) | UI/UX navigation and state management flow. |
| [ATOMIC_DESIGN_GUIDE.md](./docs/ATOMIC_DESIGN_GUIDE.md) | Guide to our frontend component organization. |
| [STORYBOOK_GUIDE.md](./docs/STORYBOOK_GUIDE.md) | How to use Storybook for component development. |
| [DESIGN_SYSTEM.md](./docs/DESIGN_SYSTEM.md) | UI design tokens and visual principles. |
| [COMPONENT_GUIDELINES.md](./docs/COMPONENT_GUIDELINES.md) | Best practices for building reusable components. |
| [REQUIREMENTS.md](./docs/REQUIREMENTS.md) | System requirements and constraints. |
| [ROADMAP.md](./docs/ROADMAP.md) | Project versioning and future milestones. |

## Event Flow & Concurrency

- **Event Driven**: Successful transactions publish `transaction.transfer.completed` to RabbitMQ, which the Notification Service consumes.
- **Concurrency**: Handled via **Optimistic Locking** in the Wallet Service, using version numbers to detect and prevent race conditions.
- **Idempotency**: Requests with the `Idempotency-Key` header are cached in Redis and persisted in Postgres to ensure exactly-once processing.
