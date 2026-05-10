# Versioning Roadmap

This document outlines the development milestones and future plans for the Distributed Payment System.

## Project Status: `v1.0.0-beta` (Current)

We have successfully established the core architecture and essential features of the payment system.

### Key Achievements
- **Microservices Architecture**: Four independent services (Auth, Wallet, Transaction, Notification) communicating via HTTP and RabbitMQ.
- **Authentication**: Robust JWT-based authentication system with Access and Refresh tokens.
- **Wallet Operations**: Implementation of balance management, top-ups, and secure transfers with row-level locking.
- **Asynchronous Notifications**: Event-driven notification system using RabbitMQ.
- **Observability Stack**: Centralized logging (Loki) and monitoring (Prometheus/Grafana) fully integrated.
- **Automated Migrations**: In-app migration logic for seamless database schema updates.
- **Frontend Dashboard**: A modern Vue 3 + Tailwind CSS dashboard providing a real-time overview of user balances and transactions.

---

## Milestone: `v1.1.0` (Upcoming)

Focus on stability, feature parity, and developer experience.

### Planned Features
- [ ] **Withdrawal Implementation**: Complete the end-to-end withdrawal flow, including UI support.
- [ ] **Comprehensive Testing**: Achieve 80%+ test coverage for core business logic in all services.
- [ ] **Enhanced Analytics**: Add interactive charts to the dashboard for income/expense trends.
- [ ] **Notification History**: Implement persistent notification storage and "Mark as Read" functionality.
- [ ] **Transaction Filtering**: Advanced search and filtering for transaction history in the UI.

---

## Milestone: `v2.0.0` (Long-term)

Scaling the system for production-grade usage and expanding capabilities.

### Future Goals
- [ ] **Kubernetes Orchestration**: Transition from Docker Compose to K8s for better scaling and resilience.
- [ ] **Multi-Currency Support**: Support for different currencies with real-time exchange rate integration.
- [ ] **Merchant Gateway**: APIs for external merchants to integrate and accept payments.
- [ ] **Audit Logging**: Comprehensive security audit logs for all administrative actions.
- [ ] **Mobile App**: Dedicated mobile client built with Capacitor or Native frameworks.

---

*Last Updated: May 8, 2026*
