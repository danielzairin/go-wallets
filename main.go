package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/danielzairin/go-wallets/wallets"
	_ "modernc.org/sqlite"
)

var DB_NAME string = "database.sqlite"

func main() {
	db, err := sql.Open("sqlite", DB_NAME)
	if err != nil {
		panic(err)
	}

	wal := wallets.NewWallet(db)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		wallets, err := wal.List()
		if err != nil {
			internalError(w, err)
			return
		}

		fmt.Fprintf(w, "wallets = %v\n", wallets)
	})

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		if name == "" {
			reject(w, http.StatusBadRequest)
			return
		}

		wallet, err := wal.Create(name)
		if err != nil {
			internalError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, wallet.Name)
	})

	http.HandleFunc("POST /top-up/{walletName}", func(w http.ResponseWriter, r *http.Request) {
		walletName := r.PathValue("walletName")

		wallet, err := wal.FindByName(walletName)
		if err != nil {
			reject(w, http.StatusNotFound)
			return
		}

		wal.TopUp(&wallet, 10)
	})

	http.HandleFunc("POST /transfer", func(w http.ResponseWriter, r *http.Request) {
		from := r.FormValue("from")
		to := r.FormValue("to")
		amount := r.FormValue("amount")

		amountNum, err := strconv.Atoi(amount)
		if err != nil {
			reject(w, http.StatusBadRequest)
			return
		}

		wal1, err := wal.FindByName(from)
		if err != nil {
			reject(w, http.StatusNotFound)
			return
		}

		wal2, err := wal.FindByName(to)
		if err != nil {
			reject(w, http.StatusNotFound)
			return
		}

		wal.Transfer(&wal1, &wal2, amountNum)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func reject(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	fmt.Fprintln(w, http.StatusText(code))
}

func internalError(w http.ResponseWriter, err error) {
	fmt.Printf("Internal Server Error: %s\n", err)
	reject(w, http.StatusInternalServerError)
}
