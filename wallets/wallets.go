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
	_, err := w.db.Exec("UPDATE wallets SET amount = amount + ? WHERE name = ?", amount, wallet.Name)
	if err != nil {
		return 0, nil
	}

	return wallet.Amount, nil
}

func (w *Wallets) List() ([]Wallet, error) {
	rows, err := w.db.Query(`SELECT name, amount FROM wallets`)
	if err != nil {
		return nil, err
	}

	wallets := make([]Wallet, 0)

	defer rows.Close()
	for rows.Next() {
		var wallet Wallet
		err = rows.Scan(&wallet.Name, &wallet.Amount)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func (w *Wallets) FindByName(name string) (Wallet, error) {
	row := w.db.QueryRow(`SELECT name, amount FROM wallets WHERE name = ?`, name)

	var wallet Wallet
	err := row.Scan(&wallet.Name, &wallet.Amount)
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}
