# notification-service internal layout

- `domain`: entities (`Notification`, `DeliveryAttempt`)
- `usecase`: event consume, notify, retry
- `repository`: notification persistence interfaces
- `delivery/http`: status endpoints
- `infrastructure`: postgres repository, rabbitmq consumer
