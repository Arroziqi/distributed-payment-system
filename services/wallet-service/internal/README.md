# wallet-service internal layout

- `domain`: entities (`Wallet`, `LedgerEntry`)
- `usecase`: topup, withdrawal, balance retrieval
- `repository`: wallet and ledger repository interfaces
- `delivery/http`: transport handlers and request validation
- `infrastructure`: postgres, redis lock, rabbitmq publisher
