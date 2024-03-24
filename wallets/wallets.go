package wallets

import (
	"database/sql"
	"fmt"
)

type Wallet struct {
	Name   string
	Amount int
}

type Wallets struct {
	db *sql.DB
}

func NewWallet(db *sql.DB) *Wallets {
	return &Wallets{db}
}

func (w *Wallets) Create(name string) (*Wallet, error) {
	if len(name) < 4 {
		return nil, fmt.Errorf("name must be at least 4 characters long")
	}

	wallet := &Wallet{Name: name, Amount: 0}

	_, err := w.db.Exec(`INSERT INTO wallets (name, amount) VALUES (?, ?)`, wallet.Name, wallet.Amount)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (w *Wallets) TopUp(wallet *Wallet, amount int) (int, error) {
	if amount <= 0 {
		return 0, fmt.Errorf("amount to top up must be a positive number")
	}

	wallet.Amount += amount

	return wallet.Amount, nil
}
