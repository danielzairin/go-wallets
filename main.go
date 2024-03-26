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
	setup()

	db, err := sql.Open("sqlite", DB_NAME)
	if err != nil {
		panic(err)
	}

	wal := wallets.NewWallet(db)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		wallets, err := wal.List()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
			return
		}

		fmt.Fprintf(w, "wallets = %v\n", wallets)
	})

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
			return
		}

		wallet, err := wal.Create(name)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, wallet.Name)
	})

	http.HandleFunc("POST /top-up/{walletName}", func(w http.ResponseWriter, r *http.Request) {
		walletName := r.PathValue("walletName")

		wallet, err := wal.FindByName(walletName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "failed to find the wallet '%s'\n", walletName)
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "failed to parse 'amount' as integer")
			return
		}

		wal1, err := wal.FindByName(from)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "wallet with name '%v' not found\n", from)
			return
		}

		wal2, err := wal.FindByName(to)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "wallet with name '%v' not found\n", to)
			return
		}

		wal.Transfer(&wal1, &wal2, amountNum)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setup() {
	db, err := sql.Open("sqlite", DB_NAME)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`DROP TABLE IF EXISTS wallets`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE wallets (name TEXT PRIMARY KEY, amount INTEGER)`)
	if err != nil {
		panic(err)
	}
}
