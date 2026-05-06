# API Usage Guide

This guide provides a step-by-step walkthrough for testing the Distributed Payment System APIs.

## Prerequisites
- All services and databases must be running via Docker Compose (`docker-compose up -d`).
- You can access Swagger UI for each service to test endpoints interactively.

## Step 1: Register a New User

1. Navigate to **Auth Service Swagger UI**: `http://localhost:8081/swagger/index.html`
2. Expand the `POST /auth/register` endpoint.
3. Click **Try it out** and enter the following payload:
   ```json
   {
     "email": "user@example.com",
     "password": "securepassword123"
   }
   ```
4. Execute the request. You should receive a `201 Created` response containing the `user_id`. Save this `user_id` for later.

## Step 2: Login and Get Token

1. Still in the Auth Service Swagger UI, expand the `POST /auth/login` endpoint.
2. Click **Try it out** and enter your credentials:
   ```json
   {
     "email": "user@example.com",
     "password": "securepassword123"
   }
   ```
3. Execute the request. You will receive an `access_token` and a `refresh_token`.
4. Copy the `access_token`.

## Step 3: Authorize in Swagger

For all protected endpoints (Wallet, Transaction, and Notification services), you must authorize using the `access_token`.

1. Open the Swagger UI for the service you want to interact with (e.g., Wallet Service at `http://localhost:8082/swagger/index.html`).
2. Click the **Authorize** button at the top right.
3. In the value field, enter: `Bearer <your_access_token>`
4. Click **Authorize** and then **Close**.

## Step 4: Create a Wallet

1. Navigate to **Wallet Service Swagger UI**: `http://localhost:8082/swagger/index.html`
2. Authorize using the steps in Step 3.
3. Expand the `POST /wallets` endpoint.
4. Provide the payload:
   ```json
   {
     "user_id": "<your_user_id_from_step_1>",
     "currency": "USD"
   }
   ```
5. Execute. You will receive a `201 Created` response with your wallet details.

*Note: You may want to register a second user and create a wallet for them to test transfers later.*

## Step 5: Top Up Wallet

1. Still in the Wallet Service Swagger UI, expand the `POST /wallet/topups` endpoint.
2. Provide the payload:
   ```json
   {
     "user_id": "<your_user_id>",
     "amount": 1000
   }
   ```
3. Execute. Your balance will be credited.

## Step 6: Transfer Money

1. Navigate to **Wallet Service Swagger UI** or **Transaction Service Swagger UI**.
   - If using Wallet Service, use `POST /wallet/transfers`.
   - If using Transaction Service (`http://localhost:8083/swagger/index.html`), use `POST /transactions/payments` and provide the Idempotency-Key header.
2. Example using Wallet Service:
   ```json
   {
     "from_user_id": "<your_user_id>",
     "to_user_id": "<recipient_user_id>",
     "amount": 250
   }
   ```
3. Execute. The funds will be transferred and the new balances will be returned.

## Step 7: Check Transaction History

1. Navigate to **Transaction Service Swagger UI**: `http://localhost:8083/swagger/index.html`
2. Authorize with your token.
3. Expand the `GET /transactions` endpoint.
4. Execute to see a paginated list of all transactions. You can also use `GET /transactions/{id}` for specific details.

## Step 8: Check Notification Logs

Notifications are generated asynchronously via RabbitMQ upon successful transactions.

1. Navigate to **Notification Service Swagger UI**: `http://localhost:8084/swagger/index.html`
2. Authorize with your token.
3. Since notifications are processed in the background and stored/sent, you can check system logs or use the `GET /notifications/{id}` endpoint if you have a specific notification ID, to verify the notification was created.
