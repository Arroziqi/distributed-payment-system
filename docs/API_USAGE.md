# API Usage

This document outlines the API endpoints consumed by the frontend.

## Common Headers
All protected routes require the `Authorization` header:
`Authorization: Bearer <jwt_token>`

## Auth Service (Port: 8081)
- `POST /api/v1/auth/register`: Register a new user. Body: `{ username, password }`
- `POST /api/v1/auth/login`: Authenticate user and receive JWT. Body: `{ username, password }`

## Wallet Service (Port: 8082)
- `POST /api/v1/wallet`: Create a new wallet for the authenticated user.
- `GET /api/v1/wallet`: Retrieve wallet details and current balance.
- `POST /api/v1/wallet/topup`: Top up the wallet balance. Body: `{ amount }`

## Transaction Service (Port: 8083)
- `POST /api/v1/transaction/transfer`: Transfer funds to another wallet. Body: `{ receiver_id, amount, description }`
- `POST /api/v1/transaction/withdraw`: Withdraw funds. Body: `{ amount, description }`
- `GET /api/v1/transaction/history`: Retrieve the transaction history for the user.

## Notification Service (Port: 8084)
- `GET /api/v1/notification`: Retrieve all recent notifications for the authenticated user.
