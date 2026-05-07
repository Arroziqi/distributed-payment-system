# System Overview

## What is this project?
The Distributed Payment System is a production-style Golang microservices application demonstrating a complete payment ecosystem. It includes authentication, wallet management, transaction processing, and asynchronous notifications, all tied together by a Quasar Framework frontend.

## Why it exists
This project serves as a showcase of advanced backend engineering and frontend architecture principles, designed to demonstrate proficiency in building scalable, observable, and resilient distributed systems to recruiters and technical peers.

## Architecture Goals
- **Separation of Concerns:** Each microservice handles a specific domain.
- **Resilience:** Services can fail independently without taking down the entire system.
- **Observability:** Centralized logging and metrics to monitor system health and trace requests.

## Distributed Systems Concepts Used
- **Microservices Architecture:** Independently deployable services.
- **Event-Driven Architecture:** Asynchronous communication via RabbitMQ.
- **Centralized Authentication:** JWT-based stateless authentication.
- **Concurrency Control:** Go routines for parallel processing and database transaction isolation levels.

## Observability Stack
- **Prometheus:** Metrics collection.
- **Grafana:** Data visualization and dashboards.
- **Loki & Promtail:** Centralized log aggregation.

## Authentication Flow
Stateless authentication using JSON Web Tokens (JWT). The Auth Service issues tokens, which are passed via the `Authorization: Bearer <token>` header and validated by API Gateways or intermediate services.

## Event-Driven Flow (RabbitMQ)
The Transaction Service publishes events (e.g., `transaction.success`) to a RabbitMQ exchange. The Notification Service consumes these events and processes them asynchronously.

## Docker Orchestration
The entire ecosystem, including the frontend, backend services, databases, message broker, and observability tools, is orchestrated using Docker Compose for seamless local development and deployment.
