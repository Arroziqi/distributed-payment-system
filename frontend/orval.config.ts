import { defineConfig } from 'orval';

export default defineConfig({
  auth: {
    input: '../services/auth-service/docs/swagger.json',
    output: {
      target: './src/api/generated/auth.ts',
      client: 'vue-query',
      httpClient: 'axios',
      mode: 'tags-split',
      override: {
        mutator: {
          path: './src/api/client.ts',
          name: 'customInstance',
        },
      },
    },
  },
  wallet: {
    input: '../services/wallet-service/docs/swagger.json',
    output: {
      target: './src/api/generated/wallet.ts',
      client: 'vue-query',
      httpClient: 'axios',
      mode: 'tags-split',
      override: {
        mutator: {
          path: './src/api/client.ts',
          name: 'customInstance',
        },
      },
    },
  },
  transaction: {
    input: '../services/transaction-service/docs/swagger.json',
    output: {
      target: './src/api/generated/transaction.ts',
      client: 'vue-query',
      httpClient: 'axios',
      mode: 'tags-split',
      override: {
        mutator: {
          path: './src/api/client.ts',
          name: 'customInstance',
        },
      },
    },
  },
  notification: {
    input: '../services/notification-service/docs/swagger.json',
    output: {
      target: './src/api/generated/notification.ts',
      client: 'vue-query',
      httpClient: 'axios',
      mode: 'tags-split',
      override: {
        mutator: {
          path: './src/api/client.ts',
          name: 'customInstance',
        },
      },
    },
  },
});
