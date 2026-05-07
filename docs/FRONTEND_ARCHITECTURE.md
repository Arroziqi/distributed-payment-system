# Frontend Architecture

## Stack
- **Framework:** Vue.js 3
- **UI Toolkit:** Quasar Framework
- **State Management:** Pinia
- **HTTP Client:** Axios
- **Component Explorer:** Storybook
- **Language:** TypeScript

## Architecture Pattern
The application strictly follows the **Atomic Design** methodology to ensure components are modular, reusable, and easily testable.

## Directory Structure
- `src/components/`: Houses all Atomic Design components.
- `src/stores/`: Pinia stores for global state management (Auth, Wallet, Transactions, Notifications).
- `src/services/`: API abstractions and Axios interceptor configurations.
- `src/styles/`: Global stylesheets and Quasar variable overrides.
- `src/router/`: Vue Router configurations and navigation guards.

## State Management Flow
1. Components trigger actions in Pinia stores.
2. Store actions invoke API services.
3. API services perform HTTP requests.
4. Responses update the store state.
5. Reactivity updates the components seamlessly.
