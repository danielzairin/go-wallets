Wallet core logic
- Wallet struct
- wallets.TopUp(to, amount)
- wallets.MakeTransaction(from, to, amount)
- wallets.List()
- wallets.FindOne()
- wallets.Create()

Setup routes
- POST /wallets
- GET /wallets
- GET /wallets/{id}
- POST /wallets/{id}/top-up
- DELETE /wallets/{id}