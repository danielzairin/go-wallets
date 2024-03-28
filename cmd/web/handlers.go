package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app App) listWallets(w http.ResponseWriter, r *http.Request) {
	wallets, err := app.wallets.List()

	if err != nil {
		internalError(w, err)
		return
	}

	fmt.Fprintf(w, "wallets = %v\n", wallets)
}

func (app App) postCreateWallet(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		reject(w, http.StatusBadRequest)
		return
	}

	wallet, err := app.wallets.Create(name)
	if err != nil {
		internalError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, wallet.Name)
}

func (app App) postTopUpWallet(w http.ResponseWriter, r *http.Request) {
	walletName := r.PathValue("walletName")

	wallet, err := app.wallets.FindByName(walletName)
	if err != nil {
		reject(w, http.StatusNotFound)
		return
	}

	app.wallets.TopUp(&wallet, 10)
}

func (app App) postTransferFunds(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	to := r.FormValue("to")
	amount := r.FormValue("amount")

	amountNum, err := strconv.Atoi(amount)
	if err != nil {
		reject(w, http.StatusBadRequest)
		return
	}

	wal1, err := app.wallets.FindByName(from)
	if err != nil {
		reject(w, http.StatusNotFound)
		return
	}

	wal2, err := app.wallets.FindByName(to)
	if err != nil {
		reject(w, http.StatusNotFound)
		return
	}

	app.wallets.Transfer(&wal1, &wal2, amountNum)
}
