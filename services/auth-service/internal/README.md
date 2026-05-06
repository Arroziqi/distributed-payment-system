# auth-service internal layout

- `domain`: entities (`User`, `RefreshToken`)
- `usecase`: register/login/refresh/logout business flows
- `repository`: interfaces for persistence and token storage
- `delivery/http`: handlers + DTO mapping
- `infrastructure`: postgres, redis, rabbitmq adapters
