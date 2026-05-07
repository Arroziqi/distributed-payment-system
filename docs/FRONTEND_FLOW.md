# Frontend UI Workflow

This document outlines the standard user flow through the frontend application.

## 1. Initial Load & Authentication Check
- The user accesses the application.
- `Vue Router` global navigation guard checks if the route requires authentication (`meta.requiresAuth`).
- It checks the `authStore` for a valid JWT.
- If unauthenticated and trying to access a protected route, the user is redirected to `/login`.

## 2. Authentication Flow
- The user inputs credentials into the `LoginForm` (Molecule).
- The form emits a `submit` event to the `LoginPage`.
- `LoginPage` calls `AuthService.login()`.
- Upon success, the JWT is stored in `localStorage` and `authStore`.
- The user is redirected to `/dashboard`.

## 3. API Interception
- All subsequent Axios requests pass through an interceptor.
- The interceptor automatically attaches the `Authorization: Bearer <token>` header.
- If any request returns a `401 Unauthorized` (e.g., token expired), the response interceptor automatically triggers `authStore.logout()` and redirects to `/login`.

## 4. Dashboard & Data Loading
- Upon entering `/dashboard`, `onMounted` hooks trigger data fetching.
- `WalletStore.fetchBalance()` and `NotificationStore.fetchNotifications()` are called.
- While data is fetching, specific component areas show `BaseLoader` or Quasar spinners.

## 5. Wallet Operations
- Users navigate to `/wallet`.
- They can input an amount and submit the Top Up form.
- The `WalletPage` calls `WalletStore.topUp()`, which hits the API and then refreshes the local balance state.
- Notifications are displayed via Quasar's `Notify` plugin indicating success or failure.
