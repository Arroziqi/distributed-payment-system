# System Flow & Business Logic

This document explains the business logic flows and system architecture sequences for the Distributed Payment System.

## 1. User Registration & Login Flow

The Auth Service handles user identity and issues JWT tokens for subsequent operations.

```mermaid
sequenceDiagram
    participant User
    participant Auth Service
    participant PostgreSQL Database
    
    %% Registration Flow
    Note over User, PostgreSQL Database: User Registration
    User->>Auth Service: POST /register (email, password)
    Auth Service->>Auth Service: Hash password
    Auth Service->>PostgreSQL Database: Insert new user
    PostgreSQL Database-->>Auth Service: Return user ID
    Auth Service-->>User: 201 Created (User ID)
    
    %% Login Flow
    Note over User, PostgreSQL Database: User Login
    User->>Auth Service: POST /login (email, password)
    Auth Service->>PostgreSQL Database: Fetch user & password hash
    PostgreSQL Database-->>Auth Service: User details
    Auth Service->>Auth Service: Validate password
    Auth Service->>Auth Service: Generate Access & Refresh Tokens
    Auth Service-->>User: 200 OK (Tokens)
```

## 2. Wallet Creation & Top Up Flow

The Wallet Service manages user balances and ensures thread-safe operations on financial records.

```mermaid
sequenceDiagram
    participant User
    participant Wallet Service
    participant PostgreSQL Database
    
    %% Top Up Flow
    Note over User, PostgreSQL Database: Wallet Top Up
    User->>Wallet Service: POST /wallet/topups (userID, amount)
    Wallet Service->>PostgreSQL Database: BEGIN Transaction
    Wallet Service->>PostgreSQL Database: SELECT FOR UPDATE (Lock Wallet)
    PostgreSQL Database-->>Wallet Service: Current Balance
    Wallet Service->>Wallet Service: Calculate New Balance (Balance + Amount)
    Wallet Service->>PostgreSQL Database: UPDATE Wallet Balance
    Wallet Service->>PostgreSQL Database: COMMIT Transaction
    Wallet Service-->>User: 200 OK (New Balance)
```

## 3. Transfer Flow

The Transaction Service orchestrates the transfer process between two wallets, enforcing idempotency and interacting asynchronously with the Notification Service.

```mermaid
sequenceDiagram
    participant User
    participant Transaction Service
    participant Wallet Service
    participant Database (Tx)
    participant RabbitMQ
    participant Notification Service
    
    User->>Transaction Service: POST /transactions/payments
    Transaction Service->>Transaction Service: Check Idempotency Key
    Transaction Service->>Database (Tx): Save Transaction (Pending)
    Transaction Service->>Wallet Service: POST /wallet/transfers (Saga/Call)
    Wallet Service->>Wallet Service: Verify balances & Lock rows
    Wallet Service->>Wallet Service: Update sender & receiver balances
    Wallet Service-->>Transaction Service: 200 OK
    Transaction Service->>Database (Tx): Update Transaction (Success)
    Transaction Service->>RabbitMQ: Publish Transfer Event
    RabbitMQ-->>Notification Service: Consume Transfer Event
    Notification Service->>Notification Service: Process & Save Notification
    Transaction Service-->>User: 200 OK (Transfer details)
```

## 4. Withdraw Flow

Withdrawing funds decreases the wallet balance, ensuring that the account does not drop below zero.

```mermaid
sequenceDiagram
    participant User
    participant Wallet Service
    participant PostgreSQL Database
    
    User->>Wallet Service: POST /wallet/withdrawals (userID, amount)
    Wallet Service->>PostgreSQL Database: BEGIN Transaction
    Wallet Service->>PostgreSQL Database: SELECT FOR UPDATE (Lock Wallet)
    PostgreSQL Database-->>Wallet Service: Current Balance
    
    alt Insufficient Balance
        Wallet Service->>PostgreSQL Database: ROLLBACK
        Wallet Service-->>User: 422 Unprocessable Entity
    else Sufficient Balance
        Wallet Service->>Wallet Service: Calculate New Balance (Balance - Amount)
        Wallet Service->>PostgreSQL Database: UPDATE Wallet Balance
        Wallet Service->>PostgreSQL Database: COMMIT Transaction
        Wallet Service-->>User: 200 OK (New Balance)
    end
```

## Concurrency Handling

The system addresses concurrency using the following mechanisms:
- **Database Locks**: Critical paths in the Wallet Service use `SELECT ... FOR UPDATE` row-level locks within transactions. This ensures that concurrent reads and writes to a single wallet's balance are serialized, preventing race conditions and race anomalies (e.g., duplicate spending).
- **Idempotency**: The Transaction Service requires an `Idempotency-Key` header for payment processing. This ensures that retried network requests do not result in duplicate transactions.
- **Asynchronous Processing**: Non-critical paths, such as generating alerts and emails, are offloaded to RabbitMQ and handled by the Notification Service, ensuring that the primary transaction pathway remains fast and unaffected by slow external integrations.
