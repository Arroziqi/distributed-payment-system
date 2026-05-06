# transaction-service internal layout

- `domain`: entities (`Transaction`, `IdempotencyKey`)
- `usecase`: transfer orchestration, history query
- `repository`: transaction and idempotency repository interfaces
- `delivery/http`: endpoint handlers and response assembly
- `infrastructure`: postgres, redis, rabbitmq
